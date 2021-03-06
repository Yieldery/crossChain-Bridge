package eth

import (
	"time"

	"github.com/fsn-dev/crossChain-Bridge/common"
	"github.com/fsn-dev/crossChain-Bridge/log"
	"github.com/fsn-dev/crossChain-Bridge/tokens/tools"
	"github.com/fsn-dev/crossChain-Bridge/types"
)

var (
	maxScanHeight          = uint64(15000)
	retryIntervalInScanJob = 3 * time.Second
	restIntervalInScanJob  = 3 * time.Second
)

// StartSwapHistoryScanJob scan job
func (b *Bridge) StartSwapHistoryScanJob() {
	if b.TokenConfig.ContractAddress == "" {
		return
	}
	log.Info("[swaphistory] start scan swap history job", "isSrc", b.IsSrc)

	isProcessed := func(txid string) bool {
		if b.IsSrc {
			return tools.IsSwapinExist(txid)
		}
		return tools.IsSwapoutExist(txid)
	}

	go b.scanFirstLoop(isProcessed)

	b.scanTransactionHistory(isProcessed)
}

func (b *Bridge) getSwapLogs(blockHeight uint64) ([]*types.RPCLog, error) {
	token := b.TokenConfig
	contractAddress := token.ContractAddress
	var logTopic string
	if b.IsSrc {
		logTopic = common.ToHex(getLogSwapinTopic())
	} else {
		logTopic = common.ToHex(getLogSwapoutTopic())
	}
	return b.GetContractLogs(contractAddress, logTopic, blockHeight)
}

func (b *Bridge) scanFirstLoop(isProcessed func(string) bool) {
	// first loop process all tx history no matter whether processed before
	log.Info("[scanhistory] start first scan loop", "isSrc", b.IsSrc)
	initialHeight := b.TokenConfig.InitialHeight
	latest := tools.LoopGetLatestBlockNumber(b)
	for height := latest; height+maxScanHeight > latest && height >= initialHeight; {
		logs, err := b.getSwapLogs(height)
		if err != nil {
			log.Error("[scanhistory] get swap logs error", "isSrc", b.IsSrc, "height", height, "err", err)
			time.Sleep(retryIntervalInScanJob)
			continue
		}
		for _, log := range logs {
			txid := log.TxHash.String()
			if !isProcessed(txid) {
				b.processTransaction(txid)
			}
		}
		if height > 0 {
			height--
		} else {
			break
		}
	}

	log.Info("[scanhistory] finish first scan loop", "isSrc", b.IsSrc)
}

func (b *Bridge) scanTransactionHistory(isProcessed func(string) bool) {
	log.Info("[scanhistory] start scan swap history loop")
	var (
		height        uint64
		rescan        = true
		initialHeight = b.TokenConfig.InitialHeight
	)
	for {
		if rescan || height < initialHeight || height == 0 {
			height = tools.LoopGetLatestBlockNumber(b)
		}
		logs, err := b.getSwapLogs(height)
		if err != nil {
			log.Error("[swaphistory] get swap logs error", "isSrc", b.IsSrc, "height", height, "err", err)
			time.Sleep(retryIntervalInScanJob)
			continue
		}
		log.Info("[scanhistory] scan swap history", "isSrc", b.IsSrc, "height", height, "count", len(logs))
		for _, log := range logs {
			txid := log.TxHash.String()
			if isProcessed(txid) {
				rescan = true
				break // rescan if already processed
			}
			b.processTransaction(txid)
		}
		if rescan {
			time.Sleep(restIntervalInScanJob)
		} else if height > 0 {
			height--
		}
	}
}
