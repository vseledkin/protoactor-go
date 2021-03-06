package routing

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/mock"
)

func TestRouterActor_Receive_AddRoute(t *testing.T) {
	state := new(testRouterState)

	a := routerActor{state: state}

	p1 := actor.NewLocalPID("p1")
	c := new(mockContext)
	c.On("Message").Return(&AddRoutee{p1})
	c.On("Watch", p1).Once()

	state.On("GetRoutees").Return(&actor.PIDSet{})
	state.On("SetRoutees", actor.NewPIDSet(p1)).Once()

	a.Receive(c)
	mock.AssertExpectationsForObjects(t, state, c)
}

func TestRouterActor_Receive_AddRoute_NoDuplicates(t *testing.T) {
	state := new(testRouterState)

	a := routerActor{state: state}

	p1 := actor.NewLocalPID("p1")
	c := new(mockContext)
	c.On("Message").Return(&AddRoutee{p1})

	state.On("GetRoutees").Return(actor.NewPIDSet(p1))

	a.Receive(c)
	mock.AssertExpectationsForObjects(t, state, c)
}

func TestRouterActor_Receive_RemoveRoute(t *testing.T) {
	state := new(testRouterState)

	a := routerActor{state: state}

	p1 := actor.NewLocalPID("p1")
	p2 := actor.NewLocalPID("p2")
	c := new(mockContext)
	c.On("Message").Return(&RemoveRoutee{p1})
	c.On("Unwatch", p1).Once()

	state.On("GetRoutees").Return(actor.NewPIDSet(p1, p2))
	state.On("SetRoutees", actor.NewPIDSet(p2)).Once()

	a.Receive(c)
	mock.AssertExpectationsForObjects(t, state, c)
}

func TestRouterActor_Receive_BroadcastMessage(t *testing.T) {
	state := new(testRouterState)
	a := routerActor{state: state}

	p1 := actor.NewLocalPID("p1")
	p2 := actor.NewLocalPID("p2")

	child := new(mockProcess)
	child.On("SendUserMessage", mock.Anything, mock.Anything, mock.Anything).Times(2)

	actor.ProcessRegistry.Add(child, "p1")
	actor.ProcessRegistry.Add(child, "p2")

	c := new(mockContext)
	c.On("Message").Return(&BroadcastMessage{"hi"})
	c.On("Sender").Return((*actor.PID)(nil))

	state.On("GetRoutees").Return(actor.NewPIDSet(p1, p2))

	a.Receive(c)
	mock.AssertExpectationsForObjects(t, state, c, child)
}
