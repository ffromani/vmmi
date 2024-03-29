package progress

import (
	libvirt "github.com/libvirt/libvirt-go"
)

type Progress struct {
	Valid             bool   `json:"valid"`
	Percentage        uint64 `json:"percentage"`
	Iteration         uint64 `json:"iteration"`
	UserDataRemaining int64  // TODO testing helper
	info              *libvirt.DomainJobInfo
}

// no error yet - no place to report atm
func NewProgress(dom *libvirt.Domain) *Progress {
	ret := &Progress{}

	info, err := dom.GetJobStats(0)
	if err != nil {
		return ret
	}

	ret.info = info
	if !IsOngoing(ret) {
		return ret
	}

	return ret.FromDomainJobInfo(info)
}

func (p *Progress) DataRemaining() int64 {
	if p.info == nil {
		return p.UserDataRemaining
	}
	return int64(p.info.DataRemaining)
}

// no error yet - no place to report atm
func (p *Progress) FromDomainJobInfo(info *libvirt.DomainJobInfo) *Progress {
	p.info = info

	if info.MemIterationSet {
		p.Valid = true
		p.Iteration = info.MemIteration
	}

	if info.DataRemainingSet && info.DataRemainingSet {
		// ported from https://github.com/oVirt/vdsm/blob/ovirt-4.2.4/lib/vdsm/virt/migration.py#L962
		p.Valid = true
		if info.DataRemaining == 0 && info.DataTotal > 0 {
			p.Percentage = 100
		} else {
			if info.DataTotal > 0 {
				p.Percentage = 100 - 100*info.DataRemaining
			}
			if p.Percentage >= 100 {
				p.Percentage = 99
			}
		}
	}
	return p
}

func IsOngoing(p *Progress) bool {
	return p != nil && p.info != nil && p.info.OperationSet && p.info.Operation == libvirt.DOMAIN_JOB_OPERATION_MIGRATION_OUT
}
