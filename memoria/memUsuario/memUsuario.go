package memUsuario

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/memsistema"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var memoria global
var MemoriaDeUsuario []byte
var Particiones []types.Particion
var BitmapParticiones []bool
var PidAParticion map[uint32]int // Mapa para rastrear la asignación de PIDs a particiones

// Funcion para iniciar la memoria y definir las particiones
func Inicializar_Memoria_De_Usuario(logger *slog.Logger) {
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
	PidAParticion = make(map[uint32]int)
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

func AsignarPID(pid uint32, tamanio_proceso int, path string) {

	algoritmo := utils.Configs.SearchAlgorithm
	switch algoritmo {
	case "FIRST":
		FirstFit(pid, tamanio_proceso, path)
	case "BEST":
		BestFit(pid, tamanio_proceso, path)
	case "WORST":
		WorstFit(pid, tamanio_proceso, path)
	}
}

// first fit
func FirstFit(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		particiones := utils.Configs.Partitions
		for i := 0; i < len(BitmapParticiones); i++ {
			if !BitmapParticiones[i] {
				if tamanio_proceso < particiones[i] {
					PidAParticion[pid] = i
					BitmapParticiones[i] = true
					fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
					memsistema.CrearContextoPID(pid, uint32(Particiones[i].Base), uint32(Particiones[i].Limite))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
					return
				}
			}
		}
		(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
	}
}

// BestFit
func BestFit(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
			return
		} else {
			PidAParticion[pid] = pos_menor
			BitmapParticiones[pos_menor] = true
			fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_menor+1)
			memsistema.CrearContextoPID(pid, uint32(Particiones[pos_menor].Base), uint32(Particiones[pos_menor].Limite))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
	}
}

func WorstFit(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
			return
		} else {
			PidAParticion[pid] = pos_mayor
			BitmapParticiones[pos_mayor] = true
			fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_mayor+1)
			memsistema.CrearContextoPID(pid, uint32(Particiones[pos_mayor].Base), uint32(Particiones[pos_mayor].Limite))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
	}
}

// Función de compactación
func Compactar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Acá va la lógica de compactación
	}
}
