package mezclas

import (


	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {

	CrearMezcla(p domain.Mezclas) (domain.Mezclas, error)
    

}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
	return &service{r}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UN NUEVO PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) CrearMezcla(p domain.Mezclas) (domain.Mezclas, error) {
    // Crear el producto utilizando el repositorio
    productoCreado, err := s.r.CrearMezcla(p)
    if err != nil {
        return domain.Mezclas{}, err
    }
    return productoCreado, nil
}