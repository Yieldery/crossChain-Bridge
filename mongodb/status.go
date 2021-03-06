package mongodb

// -----------------------------------------------
// swap status change graph
//
// TxNotStable -> |- TxVerifyFailed -> stop
//                |- TxCanRecall -> TxToBeRecall -> |- TxRecallFailed -> retry
//                                                  |- TxProcessed (->MatchTxNotStable)
//                |- TxNotSwapped -> |- TxSwapFailed -> retry
//                                   |- TxProcessed (->MatchTxNotStable)
// -----------------------------------------------
// swap result status change graph
//
// TxWithWrongMemo -> |
// MatchTxEmpty    -> | MatchTxNotStable -> MatchTxStable
// -----------------------------------------------

// SwapStatus swap status
type SwapStatus uint16

// swap status values
const (
	TxNotStable      SwapStatus = iota // 0
	TxVerifyFailed                     // 1
	TxCanRecall                        // 2
	TxToBeRecall                       // 3
	TxRecallFailed                     // 4
	TxNotSwapped                       // 5
	TxSwapFailed                       // 6
	TxProcessed                        // 7
	MatchTxEmpty                       // 8
	MatchTxNotStable                   // 9
	MatchTxStable                      // 10
	TxWithWrongMemo                    // 11
)

func (status SwapStatus) String() string {
	switch status {
	case TxNotStable:
		return "TxNotStable"
	case TxVerifyFailed:
		return "TxVerifyFailed"
	case TxCanRecall:
		return "TxCanRecall"
	case TxToBeRecall:
		return "TxToBeRecall"
	case TxRecallFailed:
		return "TxRecallFailed"
	case TxNotSwapped:
		return "TxNotSwapped"
	case TxSwapFailed:
		return "TxSwapFailed"
	case TxProcessed:
		return "TxProcessed"
	case MatchTxEmpty:
		return "MatchTxEmpty"
	case MatchTxNotStable:
		return "MatchTxNotStable"
	case MatchTxStable:
		return "MatchTxStable"
	case TxWithWrongMemo:
		return "TxWithWrongMemo"
	default:
		panic("unknown swap status")
	}
}
