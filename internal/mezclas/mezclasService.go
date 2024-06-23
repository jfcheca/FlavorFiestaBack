package mezclas

import (
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {

	CrearMezcla(p domain.Mezclas) (domain.Mezclas, error)
    BuscarMezcla(id int) (domain.Mezclas, error)
    DeleteMezclas(id int) error
    BuscarTodasLasMezclas() ([]domain.Mezclas, error)
    

}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
	return &service{r}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UN NUEVO PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) CrearMezcla(mezcla domain.Mezclas) (domain.Mezclas, error) {
    // Aquí podrías validar los datos de la mezcla si es necesario
    mezclaCreada, err := s.r.CrearMezcla(mezcla)
    if err != nil {
        return domain.Mezclas{}, err
    }
    return mezclaCreada, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE CATEGORIA POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) BuscarMezcla(id int) (domain.Mezclas, error) {
    categoria, err := s.r.BuscarMezcla(id)
    if err != nil {
        return domain.Mezclas{}, err
    }
    return categoria, nil
}

func (s *service) DeleteMezclas(id int) error {
    err := s.r.DeleteMezclas(id)
    if err != nil {
        return err
    }
    return nil
}


func (s *service) BuscarTodasLasMezclas() ([]domain.Mezclas, error) {
    productos, err := s.r.BuscarTodasLasMezclas()
    if err != nil {
        return nil, fmt.Errorf("error buscando todas las mezclas: %w", err)
    }
    return productos, nil
}