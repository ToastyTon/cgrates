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
package config

import (
	"reflect"
	"testing"

	"github.com/cgrates/cgrates/utils"
)

func TestDiameterAgentCfgloadFromJsonCfg(t *testing.T) {
	var dacfg, expected DiameterAgentCfg
	if err := dacfg.loadFromJsonCfg(nil, utils.INFIELD_SEP); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(dacfg, expected) {
		t.Errorf("Expected: %+v ,recived: %+v", expected, dacfg)
	}
	if err := dacfg.loadFromJsonCfg(new(DiameterAgentJsonCfg), utils.INFIELD_SEP); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(dacfg, expected) {
		t.Errorf("Expected: %+v ,recived: %+v", expected, dacfg)
	}
	cfgJSONStr := `{
"diameter_agent": {
	"enabled": false,											// enables the diameter agent: <true|false>
	"listen": "127.0.0.1:3868",									// address where to listen for diameter requests <x.y.z.y:1234>
	"dictionaries_path": "/usr/share/cgrates/diameter/dict/",	// path towards directory holding additional dictionaries to load
	"sessions_conns": ["*internal"],
	"origin_host": "CGR-DA",									// diameter Origin-Host AVP used in replies
	"origin_realm": "cgrates.org",								// diameter Origin-Realm AVP used in replies
	"vendor_id": 0,												// diameter Vendor-Id AVP used in replies
	"product_name": "CGRateS",									// diameter Product-Name AVP used in replies
	"synced_conn_requests": true,
	"request_processors": [],
},
}`
	expected = DiameterAgentCfg{
		Listen:           "127.0.0.1:3868",
		DictionariesPath: "/usr/share/cgrates/diameter/dict/",
		SessionSConns:    []string{utils.ConcatenatedKey(utils.MetaInternal, utils.MetaSessionS)},
		OriginHost:       "CGR-DA",
		OriginRealm:      "cgrates.org",
		VendorId:         0,
		ProductName:      "CGRateS",
		SyncedConnReqs:   true,
	}
	if jsnCfg, err := NewCgrJsonCfgFromBytes([]byte(cfgJSONStr)); err != nil {
		t.Error(err)
	} else if jsnDaCfg, err := jsnCfg.DiameterAgentJsonCfg(); err != nil {
		t.Error(err)
	} else if err = dacfg.loadFromJsonCfg(jsnDaCfg, utils.INFIELD_SEP); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expected, dacfg) {
		t.Errorf("Expected: %+v , recived: %+v", utils.ToJSON(expected), utils.ToJSON(dacfg))
	}
}

func TestDiameterAgentCfgAsMapInterface(t *testing.T) {
	cfgJSONStr := `{
	"diameter_agent": {
		"enabled": false,											
		"listen": "127.0.0.1:3868",									
		"dictionaries_path": "/usr/share/cgrates/diameter/dict/",	
		"sessions_conns": ["*internal"],
		"origin_host": "CGR-DA",									
		"origin_realm": "cgrates.org",								
		"vendor_id": 0,												
		"product_name": "CGRateS",									
		"synced_conn_requests": true,
		"request_processors": [],
	},
}`
	eMap := map[string]interface{}{
		utils.ASRTemplateCfg:        "",
		utils.ConcurrentRequestsCfg: -1,
		utils.DictionariesPathCfg:   "/usr/share/cgrates/diameter/dict/",
		utils.EnabledCfg:            false,
		utils.ForcedDisconnectCfg:   "*none",
		utils.ListenCfg:             "127.0.0.1:3868",
		utils.ListenNetCfg:          "tcp",
		utils.OriginHostCfg:         "CGR-DA",
		utils.OriginRealmCfg:        "cgrates.org",
		utils.ProductNameCfg:        "CGRateS",
		utils.RARTemplateCfg:        "",
		utils.SessionSConnsCfg:      []string{"*internal"},
		utils.SyncedConnReqsCfg:     true,
		utils.VendorIdCfg:           0,
		utils.RequestProcessorsCfg:  []map[string]interface{}{},
	}
	if cgrCfg, err := NewCGRConfigFromJsonStringWithDefaults(cfgJSONStr); err != nil {
		t.Error(err)
	} else if rcv := cgrCfg.diameterAgentCfg.AsMapInterface(utils.EmptyString); !reflect.DeepEqual(rcv, eMap) {
		t.Errorf("Expected %+v \n, received %+v", eMap, rcv)
	}
}

func TestDiameterAgentCfgAsMapInterface1(t *testing.T) {
	cfgJSONStr := `{
	"diameter_agent": {
		"enabled": true,
		"dictionaries_path": "/usr/share/cgrates/diameter",			
		"synced_conn_requests": false,
	},
}`
	eMap := map[string]interface{}{
		utils.ASRTemplateCfg:        "",
		utils.ConcurrentRequestsCfg: -1,
		utils.DictionariesPathCfg:   "/usr/share/cgrates/diameter",
		utils.EnabledCfg:            true,
		utils.ForcedDisconnectCfg:   "*none",
		utils.ListenCfg:             "127.0.0.1:3868",
		utils.ListenNetCfg:          "tcp",
		utils.OriginHostCfg:         "CGR-DA",
		utils.OriginRealmCfg:        "cgrates.org",
		utils.ProductNameCfg:        "CGRateS",
		utils.RARTemplateCfg:        "",
		utils.SessionSConnsCfg:      []string{"*internal"},
		utils.SyncedConnReqsCfg:     false,
		utils.VendorIdCfg:           0,
		utils.RequestProcessorsCfg:  []map[string]interface{}{},
	}
	if cgrCfg, err := NewCGRConfigFromJsonStringWithDefaults(cfgJSONStr); err != nil {
		t.Error(err)
	} else if rcv := cgrCfg.diameterAgentCfg.AsMapInterface(cgrCfg.generalCfg.RSRSep); !reflect.DeepEqual(rcv, eMap) {
		t.Errorf("Expected %+v \n, received %+v", eMap, rcv)
	}
}
