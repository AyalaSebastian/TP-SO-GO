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

type Particion struct {
	Registros RegCPU `json:"registros"`
}

type CPU struct {
	Contexto          ContextoEjecucion `json:"contexto"`
	MMU               MMU               `json:"mmu"`
	Memoria           Memoria           `json:"memoria"`
	InstruccionActual string            `json:"instruccion_actual"`
}
type MMU struct {
}
type Memoria struct {
}
