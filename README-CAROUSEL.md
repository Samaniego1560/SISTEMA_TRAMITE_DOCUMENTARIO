# Sistema de Gestión de Carrusel - DBU

## Descripción

Sistema completo de gestión de carrusel para la Dirección de Bienestar Universitario (DBU) que permite a los administradores cargar, configurar y gestionar imágenes del carrusel de bienvenida con texto opcional y botones de llamada a la acción.

## Características

- ✅ **Gestión Completa**: Crear, editar, eliminar carruseles
- ✅ **Upload de Imágenes**: Soporte para JPG, PNG, WEBP (máx 5MB)
- ✅ **Preview en Tiempo Real**: Vista previa antes de guardar
- ✅ **Texto Personalizable**: Título y descripción opcionales
- ✅ **Call-to-Action**: Botón con enlace personalizable
- ✅ **Control de Visibilidad**: Habilitar/deshabilitar carruseles
- ✅ **Ordenamiento**: Control del orden de visualización
- ✅ **Responsive**: Adaptado para móviles, tablets y desktop
- ✅ **Accesible**: ARIA labels y navegación por teclado

## Instalación

### Backend (PHP/Laravel)

1. **Ejecutar migración**:
```bash
cd backend-php
php artisan migrate
```

2. **Crear directorio de imágenes**:
```bash
mkdir -p public/carousel
chmod 755 public/carousel
```

3. **Verificar configuración PHP**:
Asegurar que en `php.ini`:
```ini
upload_max_filesize = 5M
post_max_size = 5M
```

### Frontend (Angular)

No requiere configuración adicional. El componente se carga automáticamente.

## Uso

### Panel de Administración

1. **Acceder al panel**:
   - Iniciar sesión como administrador (role 1)
   - Ir a menú: **CONFIGURACIÓN DE BIENVENIDA** > **CARRUSEL**

2. **Crear nuevo carrusel**:
   - Seleccionar imagen (requerido)
   - Ingresar título (opcional)
   - Ingresar descripción (opcional)
   - Ingresar texto del botón (opcional)
   - Ingresar enlace del botón (opcional)
   - Configurar orden
   - Marcar como habilitado
   - Click en "Guardar"

3. **Editar carrusel existente**:
   - Click en "Editar" en la tarjeta del carrusel
   - Modificar campos deseados
   - Click en "Actualizar"

4. **Eliminar carrusel**:
   - Click en "Eliminar" en la tarjeta del carrusel
   - Confirmar eliminación

5. **Habilitar/Deshabilitar**:
   - Click en "Activar" o "Desactivar" según el estado actual

### Visualización Pública

El carrusel se muestra automáticamente al cargar la página principal:
- Solo se muestran carruseles habilitados
- Auto-avance cada 3 segundos
- Navegación manual con flechas o indicadores
- Se muestra una vez por sesión (configurable)

## Estructura de Archivos

### Backend
```
backend-php/
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   └── CarouselSettingController.php
│   │   ├── Requests/
│   │   │   └── CarouselSetting/
│   │   │       └── CarouselSettingRequest.php
│   │   └── Resources/
│   │       └── CarouselSetting/
│   │           └── CarouselSettingResource.php
│   ├── Models/
│   │   └── CarouselSetting.php
│   └── Services/
│       └── CarouselSetting/
│           ├── CreateCarouselSettingService.php
│           ├── UpdateCarouselSettingService.php
│           ├── ListCarouselSettingService.php
│           └── DeleteCarouselSettingService.php
├── database/
│   └── migrations/
│       └── 2026_01_27_094214_create_carousel_settings_table.php
├── public/
│   └── carousel/          # Directorio de imágenes
└── routes/
    └── web.php            # Rutas API
```

### Frontend
```
front/src/app/
├── core/
│   ├── components/
│   │   └── modal-overlay/
│   │       ├── modal-overlay.component.ts
│   │       ├── modal-overlay.component.html
│   │       └── modal-overlay.component.scss
│   ├── models/
│   │   └── carousel-setting.ts
│   ├── services/
│   │   └── carousel-setting/
│   │       └── carousel-setting.service.ts
│   └── ui/
│       └── menu/
│           └── menu.component.html
└── modules/
    └── home/
        ├── pages/
        │   └── carousel-setting/
        │       ├── carousel-setting.component.ts
        │       ├── carousel-setting.component.html
        │       └── carousel-setting.component.scss
        └── home.routes.ts
```

## API Endpoints

### Públicos
- `GET /carousel/settings/public` - Obtener carruseles habilitados

### Autenticados (Admin)
- `GET /carousel/settings` - Listar todos los carruseles
- `POST /carousel/settings/create` - Crear nuevo carrusel
- `POST /carousel/settings/update/{id}` - Actualizar carrusel
- `DELETE /carousel/settings/destroy/{id}` - Eliminar carrusel

## Base de Datos

### Tabla: `carousel_settings`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| id | bigint | ID autoincremental |
| image_path | string | Nombre del archivo de imagen |
| title | string (nullable) | Título del carrusel |
| description | text (nullable) | Descripción del carrusel |
| button_text | string (nullable) | Texto del botón CTA |
| button_link | string (nullable) | URL del enlace |
| is_enabled | boolean | Estado activo/inactivo |
| order | integer | Orden de visualización |
| created_at | timestamp | Fecha de creación |
| updated_at | timestamp | Fecha de actualización |

## Configuración Avanzada

### Cambiar tiempo de auto-avance

Editar `modal-overlay.component.ts`:
```typescript
private readonly CONFIG = {
    autoAdvanceInterval: 5000  // 5 segundos
};
```

### Cambiar frecuencia de visualización

Editar `modal-overlay.component.ts`:
```typescript
private readonly CONFIG = {
    showMode: 'days',      // 'session' o 'days'
    daysUntilReshow: 7     // Días antes de volver a mostrar
};
```

### Cambiar directorio de almacenamiento

Si necesitas cambiar el directorio de imágenes, editar:

**Backend** (`CreateCarouselSettingService.php`, `UpdateCarouselSettingService.php`, `DeleteCarouselSettingService.php`):
```php
$publicPath = public_path('carousel');  // Cambiar 'carousel' por tu directorio
```

**Modelo** (`CarouselSetting.php`):
```php
public function getImageUrlAttribute(): string
{
    return url('carousel/' . $this->image_path);  // Cambiar 'carousel'
}
```

## Solución de Problemas

### Las imágenes no se cargan
1. Verificar que el directorio `public/carousel/` existe
2. Verificar permisos del directorio (755)
3. Verificar que las imágenes se guardaron correctamente
4. Revisar la consola del navegador para errores

### Error al subir imágenes
1. Verificar límites de PHP (`upload_max_filesize`, `post_max_size`)
2. Verificar permisos de escritura en `public/carousel/`
3. Verificar formato de imagen (solo JPG, PNG, WEBP)
4. Verificar tamaño de imagen (máx 5MB)

### El carrusel no aparece
1. Verificar que hay al menos un carrusel habilitado
2. Limpiar sessionStorage/localStorage del navegador
3. Verificar que el componente `modal-overlay` está incluido en la página
4. Revisar la consola del navegador para errores de API

## Seguridad

- ✅ Validación de tipo de archivo en backend
- ✅ Límite de tamaño de archivo (5MB)
- ✅ Nombres de archivo únicos (timestamp + uniqid)
- ✅ Rutas protegidas con middleware de autenticación
- ✅ Solo administradores (role 1) pueden gestionar carruseles
- ✅ Sanitización de URLs en enlaces externos
- ✅ Atributo `rel="noopener noreferrer"` en enlaces externos

## Soporte

Para reportar problemas o solicitar nuevas características, contactar al equipo de desarrollo de DBU.

---

**Desarrollado para**: Dirección de Bienestar Universitario - UNAS  
**Versión**: 1.0.0  
**Fecha**: Enero 2026
