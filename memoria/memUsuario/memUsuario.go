package memUsuario

import (
	"fmt"
	"net/http"

	memsistema "github.com/sisoputnfrba/tp-golang/memoria/memSistema"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var memoria global
var MemoriaDeUsuario []byte
var Particiones []types.Particion
var BitmapParticiones []bool
var BitMapParticionesDinamicas []bool
var PidAParticion map[uint32]int // Mapa para rastrear la asignación de PIDs a particiones

// Funcion para iniciar la memoria y definir las particiones
func Inicializar_Memoria_De_Usuario() {
	// Inicializar el espacio de memoria con 1024 bytes
	MemoriaDeUsuario = make([]byte, utils.Configs.MemorySize)

	// Asignar las particiones fijas en la memoria usando los datos de config
	var base uint32 = 0
	for i, limite := range utils.Configs.Partitions {
		particion := types.Particion{
			Base:   base,
			Limite: uint32(limite),
		}
		Particiones = append(Particiones, particion)
		fmt.Printf("Partición %d inicializada: Base = %d, Límite = %d\n", i+1, particion.Base, particion.Limite)
		base += uint32(limite)
	}
	// Inicializar el bitmap y el mapa de PIDs
	//todas las particiones estan libres = false
	BitmapParticiones = make([]bool, len(Particiones))
	BitMapParticionesDinamicas = make([]bool, 1)
	PidAParticion = make(map[uint32]int)
}

func dividirParticion() {
	var bitMapActualizado [][]bool
	for i := 0; i < len(BitMapParticionesDinamicas); i++ {
		end := i
		if !BitMapParticionesDinamicas[i] {
			end = i + 1
			bitMapActualizado = append(bitMapActualizado, BitMapParticionesDinamicas[i:end])
			i += 1
		} else {
			bitMapActualizado = append(bitMapActualizado, BitMapParticionesDinamicas[i:end])
		}
	}

}

// Función para liberar una partición por PID
func LiberarParticionPorPID(pid uint32) error {
	particion, existe := PidAParticion[pid]
	if !existe {
		return fmt.Errorf("no se encontró el proceso %d asignado a ninguna partición", pid)
	}

	// Liberar la partición y actualizar el bitmap
	BitmapParticiones[particion] = false
	delete(PidAParticion, pid) // Eliminar la entrada del mapa
	fmt.Printf("Proceso %d liberado de la partición %d\n", pid, particion+1)
	return nil
}

func AsignarPID(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var asigno = false
		algoritmo := utils.Configs.SearchAlgorithm
		esquema := utils.Configs.Scheme
		if esquema == "FIJAS" {
			switch algoritmo {
			case "FIRST":
				asigno = FirstFitFijo(pid, tamanio_proceso, path)

			case "BEST":
				asigno = BestFitFijo(pid, tamanio_proceso, path)
			case "WORST":
				asigno = WorstFitFijo(pid, tamanio_proceso, path)
			}
			if !asigno {
				(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
		} else {
			switch algoritmo {
			case "FIRST":
				asigno = FirstFitDinamico(pid, tamanio_proceso, path)
			case "BEST":
				asigno = BestFitDinamico(pid, tamanio_proceso, path)
			case "WORST":
				asigno = WorstFitDinamico(pid, tamanio_proceso, path)
			}
			if asigno {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			} else {
				//avisar a kernel para que de la orden de compactar
			}
		}
	}
}

// first fit para particiones fijas
func FirstFitFijo(pid uint32, tamanio_proceso int, path string) bool {

	particiones := utils.Configs.Partitions
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particiones[i] {
				PidAParticion[pid] = i
				BitmapParticiones[i] = true
				fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
				memsistema.CrearContextoPID(pid, uint32(Particiones[i].Base), uint32(Particiones[i].Limite))
				return true
			}
		}
	}
	return false
}

// best fit para particiones fijas
func BestFitFijo(pid uint32, tamanio_proceso int, path string) bool {

	particiones := utils.Configs.Partitions
	var menor = 0
	var pos_menor = -1
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particiones[i] {
				if particiones[i] < menor {
					menor = particiones[i]
					pos_menor = i
				}
			}
		}
	}
	if pos_menor == -1 {
		return false
	} else {
		PidAParticion[pid] = pos_menor
		BitmapParticiones[pos_menor] = true
		fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_menor+1)
		memsistema.CrearContextoPID(pid, uint32(Particiones[pos_menor].Base), uint32(Particiones[pos_menor].Limite))
		return true
	}
}

// worst fit para particiones fijas
func WorstFitFijo(pid uint32, tamanio_proceso int, path string) bool {
	particiones := utils.Configs.Partitions
	var mayor = 0
	var pos_mayor = -1
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particiones[i] {
				if particiones[i] > mayor {
					mayor = particiones[i]
					pos_mayor = i
				}
			}
		}
	}
	if pos_mayor == -1 {
		return false
	} else {
		PidAParticion[pid] = pos_mayor
		BitmapParticiones[pos_mayor] = true
		fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_mayor+1)
		memsistema.CrearContextoPID(pid, uint32(Particiones[pos_mayor].Base), uint32(Particiones[pos_mayor].Limite))
		return true
	}
}

// empiezo con un solo espacio de memoria de 1024 bytes, si no esta reservado lo hago con el pid entrante, sino no hay espacio
func FirstFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	if !BitMapParticionesDinamicas[0] {
		BitMapParticionesDinamicas[0] = true
		fmt.Printf("Proceso %d asignado a la partición %d\n", pid, 1)
		memsistema.CrearContextoPID(pid, uint32(Particiones[0].Base), uint32(Particiones[0].Limite))
		return true
	} else {
		return false
	}
}
func BestFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	for i := 0; i < len(BitMapParticionesDinamicas); i++ {
		if !BitMapParticionesDinamicas[i] {
			if tamanio_proceso < Particion[i] {
				if tamanio_proceso < particionDividida {
					dividirParticion()
					BestFitDinamico(pid, tamanio_proceso, path)
				} else {
					PidAParticion[pid] = i
					BitmapParticiones[i] = true
					fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
					memsistema.CrearContextoPID(pid, uint32(Particiones[i].Base), uint32(Particiones[i].Limite))
					return true
				}
			}
		}
	}
	return false
}

// empiezo con un solo espacio de memoria de 1024 bytes, si no esta reservado lo hago con el pid entrante, sino no hay espacio
func WorstFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	if !BitMapParticionesDinamicas[0] {
		PidAParticion[pid] = 0
		BitMapParticionesDinamicas[0] = true
		return true
	} else {
		return false
	}
}
