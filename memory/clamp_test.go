package memory

import (
	"context"
	"testing"
)

// Um limite acima do maximo deve clampar para o MAXIMO, nao para o padrao.
//
// O padrao (50) e menor que o maximo (200), entao pedir 201 devolvia menos que
// pedir 200: nao monotonico e silencioso. Isso quebra qualquer cliente que peca
// "os ultimos N eventos" para calcular o que mudou desde a ultima execucao --
// ele recebe 50, conclui que houve pouco movimento, e nao age.
func TestListEventsClampsAboveMaxToMaxNotDefault(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	const total = 60 // mais que o padrao (50), menos que o maximo (200)
	for i := 0; i < total; i++ {
		if err := s.AppendEvent(ctx, Event{
			Kind:    "memory.test",
			Payload: "clamp",
			Source:  Source{Kind: "test", Ref: "clamp"},
		}); err != nil {
			t.Fatal(err)
		}
	}

	atMax, err := s.ListEvents(ctx, 200)
	if err != nil {
		t.Fatal(err)
	}
	if len(atMax) != total {
		t.Fatalf("ListEvents(200) = %d eventos, esperado %d", len(atMax), total)
	}

	aboveMax, err := s.ListEvents(ctx, 201)
	if err != nil {
		t.Fatal(err)
	}
	if len(aboveMax) != total {
		t.Fatalf("ListEvents(201) = %d eventos, esperado %d: pedir acima do maximo "+
			"nao pode devolver menos que pedir o maximo", len(aboveMax), total)
	}
}
