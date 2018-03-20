package got

import (
	"time"
	"github.com/giuseppesalvo/tm"
)

type StateIndex int
type StateFn func( ctx *ConversationalCtx )
type States []StateFn

type UserAnswer struct {
	Answer     string
	StateIndex StateIndex
}

type Session struct {
	UserId         string
	StateIndex     StateIndex
	Answers        []UserAnswer
	CreatedAt      time.Time
	Cronology 	   []Message
	RemindInterval *tm.Interval
	ExpireTimeout  *tm.Timeout
	Data 		   interface{}
	Running    	   bool
}

func ( session *Session ) End( ctx *ConversationalCtx ) {
	ctx.Plugin.endSession(ctx.Bot, ctx.Answer, session)
}

func ( session *Session ) GoBack( ctx *ConversationalCtx ) {
	session.StateIndex -= 1
	ctx.Plugin.run(ctx.Bot, ctx.Answer)
}

func ( session *Session ) Error( ctx *ConversationalCtx ) {
	// do nothing for now
}

func ( session *Session ) StayHere( ctx *ConversationalCtx ) {
	// do nothing for now
}

func ( session *Session ) GoTo( ctx *ConversationalCtx, index StateIndex ) {
	session.StateIndex = index
	ctx.Plugin.run(ctx.Bot, ctx.Answer)
}

func ( session *Session ) GoToStart( ctx *ConversationalCtx ) {
	session.GoTo(ctx, 0)
}

func ( session *Session ) WaitForAnswer( ctx *ConversationalCtx ) {
	session.StateIndex += 1
}

func ( session *Session ) SkipToNext( ctx *ConversationalCtx ) {
	session.StateIndex += 1
	ctx.Plugin.run(ctx.Bot, ctx.Answer)
}