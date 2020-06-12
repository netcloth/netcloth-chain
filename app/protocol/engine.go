package protocol

import (
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
)

type ProtocolEngine struct {
	protocols      map[uint64]Protocol
	current        uint64
	next           uint64
	ProtocolKeeper sdk.ProtocolKeeper
}

func NewProtocolEngine(protocolKeeper sdk.ProtocolKeeper) ProtocolEngine {
	engine := ProtocolEngine{
		make(map[uint64]Protocol),
		0,
		0,
		protocolKeeper,
	}
	return engine
}

func (pe *ProtocolEngine) LoadProtocol(version uint64) {
	p, flag := pe.protocols[version]
	if flag == false {
		panic("unknown protocol version!!!")
	}
	p.LoadContext()
	pe.current = version
}

func (pe *ProtocolEngine) LoadCurrentProtocol(kvStore sdk.KVStore) (bool, uint64) {
	current := pe.ProtocolKeeper.GetCurrentVersionByStore(kvStore)
	p, flag := pe.protocols[current]
	if flag {
		p.LoadContext()
		pe.current = current
	}
	return flag, current
}

func (pe *ProtocolEngine) Activate(version uint64) bool {
	protocol, flag := pe.protocols[version]
	if flag {
		protocol.Init()
		protocol.LoadContext()
		pe.current = version
	}
	return flag
}

func (pe *ProtocolEngine) GetCurrentProtocol() Protocol {
	return pe.protocols[pe.current]
}

func (pe *ProtocolEngine) GetCurrentVersion() uint64 {
	return pe.current
}

// GetUpgradeConfigByStore gets upgrade config from store
func (pe *ProtocolEngine) GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig sdk.UpgradeConfig,
	found bool) {
	return pe.ProtocolKeeper.GetUpgradeConfigByStore(store)
}

func (pe *ProtocolEngine) Add(p Protocol) Protocol {
	if p.GetVersion() != pe.next {
		panic(fmt.Errorf("Wrong version being added to the protocol engine: %d; Expecting %d", p.GetVersion(), pe.next))
	}
	pe.protocols[pe.next] = p
	pe.next++
	return p
}

// GetProtocolKeeper gets protocol keeper from engine
func (pe *ProtocolEngine) GetProtocolKeeper() sdk.ProtocolKeeper {
	return pe.ProtocolKeeper
}

func (pe *ProtocolEngine) GetByVersion(v uint64) (Protocol, bool) {
	p, flag := pe.protocols[v]
	return p, flag
}
