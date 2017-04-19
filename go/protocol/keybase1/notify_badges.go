// Auto-generated by avdl-compiler v1.3.13 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/keybase1/notify_badges.avdl

package keybase1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type ChatConversationID []byte
type BadgeState struct {
	NewTlfs       int                     `codec:"newTlfs" json:"newTlfs"`
	RekeysNeeded  int                     `codec:"rekeysNeeded" json:"rekeysNeeded"`
	NewFollowers  int                     `codec:"newFollowers" json:"newFollowers"`
	Conversations []BadgeConversationInfo `codec:"conversations" json:"conversations"`
}

type BadgeConversationInfo struct {
	ConvID         ChatConversationID `codec:"convID" json:"convID"`
	UnreadMessages int                `codec:"UnreadMessages" json:"UnreadMessages"`
}

type BadgeStateArg struct {
	BadgeState BadgeState `codec:"badgeState" json:"badgeState"`
}

type NotifyBadgesInterface interface {
	BadgeState(context.Context, BadgeState) error
}

func NotifyBadgesProtocol(i NotifyBadgesInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.NotifyBadges",
		Methods: map[string]rpc.ServeHandlerDescription{
			"badgeState": {
				MakeArg: func() interface{} {
					ret := make([]BadgeStateArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]BadgeStateArg)
					if !ok {
						err = rpc.NewTypeError((*[]BadgeStateArg)(nil), args)
						return
					}
					err = i.BadgeState(ctx, (*typedArgs)[0].BadgeState)
					return
				},
				MethodType: rpc.MethodNotify,
			},
		},
	}
}

type NotifyBadgesClient struct {
	Cli rpc.GenericClient
}

func (c NotifyBadgesClient) BadgeState(ctx context.Context, badgeState BadgeState) (err error) {
	__arg := BadgeStateArg{BadgeState: badgeState}
	err = c.Cli.Notify(ctx, "keybase.1.NotifyBadges.badgeState", []interface{}{__arg})
	return
}
