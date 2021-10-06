// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package gnbupfworker

import (
	"fmt"
	"gnbsim/common"
	"gnbsim/gnodeb/context"
	"gnbsim/logger"
	"gnbsim/util/test"

	"github.com/free5gc/ngap/ngapType"
)

func Init(gnbUpf *context.GnbUpf) {
	if gnbUpf == nil {
		logger.GNodeBLog.Errorln("GnbUpf context is nil")
		return
	}
	for {
		msg := <-gnbUpf.ReadChan
		err := HandleMessage(gnbUpf, msg)
		if err != nil {
			gnbUpf.Log.Errorln("Gnb Upf Worker HandleMessage() returned:", err)
		}
	}
}

// HandleMessage decodes an incoming GTP-U message and routes it to the corresponding
// handlers.
func HandleMessage(gnbUpf *context.GnbUpf, msg common.InterfaceMessage) error {
	// decoding the incoming packet
	tMsg := msg.(*common.TransportMessage)
	gtpHdr := &test.GtpHdr{}
	optGtpHdr := &test.GtpHdrOpt{}
	pdu, err := test.DecodeGTPv1Header(tMsg.RawPkt, gtpHdr, optGtpHdr)
	if err != nil {
		gnbUpf.Log.Errorln("DecodeGTPv1Header() returned:", err)
		return fmt.Errorf("failed to decode gtp-u header")
	}
	switch gtpHdr.MsgType {
	case test.TYPE_GPDU:
		/* A G-PDU is T-PDU encapsulated with GTP-U header*/
		err = HandleDlGpduMessage(gnbUpf, gtpHdr, optGtpHdr, pdu)
		if err != nil {
			gnbUpf.Log.Errorln("HandleDlGpduMessage() returned:", err)
			return fmt.Errorf("failed to handle downling gpdu message")
		}

		/* TODO: Handle More GTP-PDU types eg. Error Indication */
	}

	return nil
}

func SendToGnbUe(gnbue *context.GnbCpUe, event common.EventType, ngapPdu *ngapType.NGAPPDU) {
	amfmsg := common.N2Message{}
	amfmsg.Event = event
	amfmsg.Interface = common.N2_INTERFACE
	amfmsg.NgapPdu = ngapPdu
	gnbue.ReadChan <- &amfmsg
}
