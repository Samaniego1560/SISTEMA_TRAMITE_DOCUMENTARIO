package notificacion

// Puedes definir interfaces si tuvieras más servicios, pero aquí basta con uno principal.
type ServerNotificacion struct {
	SrvNotificacion PortsServerNotificacion
}

// Define la interfaz (puerto) del servicio, si quieres un nivel de abstracción
type PortsServerNotificacion interface {
	NotificarPorFaltaID(faltaID string) error
}

// Implementa el constructor del server
func NewServerNotificacion(srv PortsServerNotificacion) *ServerNotificacion {
	return &ServerNotificacion{
		SrvNotificacion: srv,
	}
}
