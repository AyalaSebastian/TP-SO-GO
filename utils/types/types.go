package types

type HandShake struct {
	Mensaje string `json:"mensaje"`
}

// --------------------------------- KERNEL ---------------------------------
type PCB struct {
	PID    uint32            `json:"pid"`
	TCBs   map[uint32]TCB    `json:"tcb"`
	Mutexs map[string]string `json:"mutexs"` // el valor es LIBRE o tid que lo tiene
}

type TCB struct {
	TID       uint32 `json:"tid"` // EL TID TAMBIEN ES SU POSICION EN EL SLICE DE TCBs
	Prioridad int    `json:"prioridad"`
	PID       uint32 `json:"pid"` //PID del proceso al que pertenece
	Quantum   int    `json:"quantum"`
}

type PathTamanio struct {
	Path    string `json:"path"`
	Tamanio int    `json:"tamanio"`
}

type EnviarHiloAMemoria struct {
	TID  uint32 `json:"tid"`
	PID  uint32 `json:"pid"`
	Path string `json:"path"`
}

type PIDTID struct {
	TID uint32 `json:"tid"`
	PID uint32 `json:"pid"`
}

type ProcessCreateParams struct {
	Path      string `json:"path"`
	Tamanio   int    `json:"tamanio"`
	Prioridad int    `json:"prioridad"`
}

type ThreadCreateParams struct {
	Path      string `json:"path"`
	Prioridad int    `json:"prioridad"`
}

type HiloDesalojado struct {
	TID    uint32 `json:"tid"`
	PID    uint32 `json:"pid"`
	Motivo string `json:"motivo"`
}

type ProcesoNew struct {
	PCB       PCB    `json:"pcb"`
	Pseudo    string `json:"pseudo"`
	Tamanio   int    `json:"tamanio"`
	Prioridad int    `json:"prioridad"`
}

// --------------------------------- CPU ---------------------------------
type RegCPU struct {
	PC     uint32 `json:"pc"`     // Program Counter (Proxima instruccion a ejecutar)
	AX     uint32 `json:"ax"`     // Registro Numerico de proposito general
	BX     uint32 `json:"bx"`     // Registro Numerico de proposito general
	CX     uint32 `json:"cx"`     // Registro Numerico de proposito general
	DX     uint32 `json:"dx"`     // Registro Numerico de proposito general
	EX     uint32 `json:"ex"`     // Registro Numerico de proposito general
	FX     uint32 `json:"fx"`     // Registro Numerico de proposito general
	GX     uint32 `json:"gx"`     // Registro Numerico de proposito general
	HX     uint32 `json:"hx"`     // Registro Numerico de proposito general
	Base   uint32 `json:"base"`   // Direccion base de la particion del proceso
	Limite uint32 `json:"limite"` // Tamanio de la particion del proceso
}

type ContextoEjecucion struct {
	Registros RegCPU `json:"registros"`
}

/*
type Particion struct {
	Registros RegCPU `json:"registros"`
}
*/

// --------------------------------- Memoria ---------------------------------

type UpdateMemoria struct {
	PID    int    `json:"pid"`
	TID    int    `json:"tid"`
	RegCPU RegCPU `json:"regCPU"` // Nuevos valores de los registros a actualizar
}

type ContextoEjecucionPID struct {
	PID    int    `json:"pid"`
	Base   uint32 `json:"base"`
	Limite uint32 `json:"limite"`
}

type ContextoEjecucionTID struct {
	TID                int               `json:"tid"`
	PC                 uint32            `json:"pc"` // Program Counter (Proxima instruccion a ejecutar)
	AX                 uint32            `json:"ax"` // Registro Numerico de proposito general
	BX                 uint32            `json:"bx"` // Registro Numerico de proposito general
	CX                 uint32            `json:"cx"` // Registro Numerico de proposito general
	DX                 uint32            `json:"dx"` // Registro Numerico de proposito general
	EX                 uint32            `json:"ex"` // Registro Numerico de proposito general
	FX                 uint32            `json:"fx"` // Registro Numerico de proposito general
	GX                 uint32            `json:"gx"` // Registro Numerico de proposito general
	HX                 uint32            `json:"hx"` // Registro Numerico de proposito general
	LISTAINSTRUCCIONES map[string]string `json:"LISTAINSTRUCCIONES"`
}

// Estructura para representar una partición de memoria
type Particion struct {
	Base   uint32
	Limite uint32
}
