package memSistema

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Mapas para almacenar los contextos de ejecución
var ContextosPID = make(map[int]types.ContextoEjecucionPID) // Contexto por PID
var ContextosTID = make(map[int]types.ContextoEjecucionTID) // Contexto por TID

// Función para inicializar un contexto de ejecución de un proceso (PID)
func crearContextoPID(pid int, base, limite uint32) {
	ContextosPID[pid] = types.ContextoEjecucionPID{
		PID:    pid,
		Base:   base,
		Limite: limite,
	}
	fmt.Printf("Contexto PID %d inicializado con Base = %d, Límite = %d\n", pid, base, limite)
}

// Función para inicializar un contexto de ejecución de un hilo (TID)
func crearContextoTID(tid int) {
	ContextosTID[tid] = types.ContextoEjecucionTID{
		TID:                tid,
		PC:                 0,
		AX:                 0,
		BX:                 0,
		CX:                 0,
		DX:                 0,
		EX:                 0,
		FX:                 0,
		GX:                 0,
		HX:                 0,
		LISTAINSTRUCCIONES: make(map[string]string), // pseudocodigo
	}
	fmt.Printf("Contexto TID %d inicializado con registros en 0\n", tid)
}

//Funcion para cargar el archivo de pseudocodigo
func CargarPseudocodigo(pid int, tid int, path string){
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("Error al abrir el archivo %s: %v", filePath, err)
    }
    defer file.Close()

    //Si no existe para el PID TID, lo creo
    if _, exists := contextosEjecucion[PID]; !exists {
        contextosEjecucion[PID] = make(map[int]*types.ContextoEjecucionTID)
    }
    if _, exists := contextosEjecucion[PID][TID]; !exists {
        contextosEjecucion[PID][TID] = &types.ContextoEjecucionTID{
            TID:                TID,
            PC:                 0,
            AX:                 0,
            BX:                 0,
            CX:                 0,
            DX:                 0,
            EX:                 0,
            FX:                 0,
            GX:                 0,
            HX:                 0,
            LISTAINSTRUCCIONES: make(map[string]string),
        }
    }
    contexto := contextosEjecucion[PID][TID]
    scanner := bufio.NewScanner(file)
    instruccionNum := 0 // Indice de instrucciones
    
	//Empiezo a leer y guardo linea x linea
    for scanner.Scan() {
        linea := scanner.Text()
        contexto.LISTAINSTRUCCIONES[strconv.Itoa(instruccionNum)] = linea
        instruccionNum++
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("Error al leer el archivo %s: %v", filePath, err)
    }
    return nil
}

func BuscarSiguienteInstruccion(tid uint32, pc uint32) {

	contexto, existeTID := ContextosTID[tid]

	if !existeTID {
		http.Error(w, "TID no encontrado", http.StatusNotFound)
		return
	}

	indiceInstruccion := pc+1

	instruccion, existe := contexto.LISTAINSTRUCCIONES[fmt.Sprintf("instr_%d", indiceInstruccion)]
    if !existe {
        return "", errors.New(fmt.Sprintf("Instrucción no encontrada para PC %d en TID %d", pc, tid))
    }
	//Log obligatorio
	fmt.Printf("Obtener instruccion: ## Obtener instrucción - (PID:TID) - (%d:%d) - Instrucción: %s \n", tid, tid, instruccion)

    return instruccion, nil


}
