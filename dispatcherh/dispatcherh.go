/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package dispatcherh

import (
	"fmt"
	"time"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

// NewDispatcherHService constructs a DispatcherHService
func NewDispatcherHService(cfg *config.CGRConfig,
	connMgr *engine.ConnManager) *DispatcherHostsService {
	return &DispatcherHostsService{
		cfg:     cfg,
		connMgr: connMgr,
		stop:    make(chan struct{}),
	}
}

// DispatcherHostsService  is the service handling dispatching towards internal components
// designed to handle automatic partitioning and failover
type DispatcherHostsService struct {
	cfg     *config.CGRConfig
	connMgr *engine.ConnManager
	stop    chan struct{}
}

// ListenAndServe will initialize the service
func (dhS *DispatcherHostsService) ListenAndServe() {
	utils.Logger.Info("Starting DispatcherH service")
	for {
		dhS.registerHosts()
		select {
		case <-dhS.stop:
			return
		case <-time.After(dhS.cfg.DispatcherHCfg().RegisterInterval):
		}
	}
}

// Shutdown is called to shutdown the service
func (dhS *DispatcherHostsService) Shutdown() {
	utils.Logger.Info(fmt.Sprintf("<%s> service shutdown initialized", utils.DispatcherH))
	dhS.unregisterHosts()
	close(dhS.stop)
	utils.Logger.Info(fmt.Sprintf("<%s> service shutdown complete", utils.DispatcherH))
	return
}

func (dhS *DispatcherHostsService) registerHosts() {
	for _, connID := range dhS.cfg.DispatcherHCfg().DispatchersConns {
		for tnt, hostCfgs := range dhS.cfg.DispatcherHCfg().Hosts {
			if tnt == utils.MetaDefault {
				tnt = dhS.cfg.GeneralCfg().DefaultTenant
			}
			args, err := NewRegisterArgs(dhS.cfg, tnt, hostCfgs)
			if err != nil {
				continue
			}
			var rply string
			if err := dhS.connMgr.Call([]string{connID}, nil, utils.DispatcherHv1RegisterHosts, args, &rply); err != nil {
				utils.Logger.Warning(fmt.Sprintf("<%s> Unable to set the hosts to the conn with ID <%s> because : %s",
					utils.DispatcherH, connID, err))
				continue
			}
		}
	}
	return
}

func (dhS *DispatcherHostsService) unregisterHosts() {
	var rply string
	for _, connID := range dhS.cfg.DispatcherHCfg().DispatchersConns {
		for tnt, hostCfgs := range dhS.cfg.DispatcherHCfg().Hosts {
			if tnt == utils.MetaDefault {
				tnt = dhS.cfg.GeneralCfg().DefaultTenant
			}
			if err := dhS.connMgr.Call([]string{connID}, nil, utils.DispatcherHv1UnregisterHosts, NewUnregisterArgs(tnt, hostCfgs), &rply); err != nil {
				utils.Logger.Warning(fmt.Sprintf("<%s> Unable to set the hosts with tenant<%s> to the conn with ID <%s> because : %s",
					utils.DispatcherH, tnt, connID, err))
				continue
			}
		}
	}
}

// Call only to implement rpcclient.ClientConnector interface
func (*DispatcherHostsService) Call(_ string, _, _ interface{}) error {
	return utils.ErrNotImplemented
}
