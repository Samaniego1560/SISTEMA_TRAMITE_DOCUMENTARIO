# 📋 Modal Carousel Institucional - Angular 17

## 🎯 Sistema de Carrusel de Comunicados para DBU

Carrusel de imágenes institucional implementado como componente **standalone de Angular 17** para mostrar múltiples comunicados, avisos y campañas de la Dirección de Bienestar Universitario.

---

## ✨ Características Principales

✅ **Auto-avance** - Cambio automático cada 3 segundos  
✅ **Pausa inteligente** - Click en imagen para detener el auto-avance  
✅ **Navegación manual** - Flechas laterales para control manual  
✅ **Teclado** - Flechas ← → para navegar, ESC para cerrar  
✅ **Indicadores visuales** - Puntos en la parte inferior  
✅ **Responsive** - Adaptado para desktop y móvil  
✅ **Accesibilidad** - ARIA labels y navegación por teclado  

---

## 📦 Ubicación de Archivos

```
front/src/app/core/components/modal-overlay/
├── modal-overlay.component.ts        # Lógica del carrusel
├── modal-overlay.component.html      # Template del carrusel
└── modal-overlay.component.scss      # Estilos del carrusel
```

---

## 🚀 Funcionamiento Automático

El carrusel se mostrará **automáticamente** cuando:

1. El usuario cargue la aplicación
2. No lo haya visto antes (según configuración)
3. Avanzará automáticamente cada **3 segundos**
4. Se detendrá cuando el usuario haga **click en la imagen**
5. Continuará con las **flechas** o **teclas del teclado**

---

## 🎨 Configurar Tus Imágenes

### **Edita el array de imágenes**

Abre `modal-overlay.component.ts` (líneas 44-61):

```typescript
images = [
  {
    src: 'assets/img/slide-1.jpg',
    alt: 'Comunicado institucional 1 - DBU',
    fallback: 'https://via.placeholder.com/800x500/1e3a8a/ffffff?text=Slide+1'
  },
  {
    src: 'assets/img/slide-2.jpg',
    alt: 'Comunicado institucional 2 - DBU',
    fallback: 'https://via.placeholder.com/800x500/3b82f6/ffffff?text=Slide+2'
  },
  // Agrega más imágenes aquí...
];
```

### **Agregar tus propias imágenes**

1. Coloca tus imágenes en: `front/src/assets/img/`
2. Por ejemplo: `comunicado-enero.jpg`, `campaña-becas.jpg`, etc.
3. Actualiza el array:

```typescript
images = [
  {
    src: 'assets/img/comunicado-enero.jpg',
    alt: 'Comunicado de Enero 2026',
    fallback: 'URL_de_respaldo' // Opcional
  },
  {
    src: 'assets/img/campaña-becas.jpg',
    alt: 'Campaña de Becas 2026',
    fallback: 'URL_de_respaldo'
  }
];
```

---

## ⚙️ Configuración del Auto-Avance

### **Cambiar el tiempo de auto-avance**

Edita `modal-overlay.component.ts` (línea 36):

```typescript
// Tiempo en milisegundos (3000 = 3 segundos)
autoAdvanceInterval: 3000  // Cambia este valor
```

**Ejemplos:**
- `2000` = 2 segundos
- `5000` = 5 segundos
- `10000` = 10 segundos

### **Frecuencia de visualización del modal**

```typescript
private readonly CONFIG = {
  showMode: 'session',     // 'session' o 'days'
  daysUntilReshow: 7,      // Solo si showMode = 'days'
  autoAdvanceInterval: 3000
};
```

---

## 🎮 Controles de Usuario

| Acción | Método |
|--------|--------|
| **Avanzar** | Click en flecha → o tecla →  |
| **Retroceder** | Click en flecha ← o tecla ← |
| **Pausar auto-avance** | Click en la imagen |
| **Ir a slide específico** | Click en indicador (punto) |
| **Cerrar modal** | Click en X o tecla ESC |

---

## 🎨 Personalizar Estilos

### **Cambiar colores de los botones**

Edita `modal-overlay.component.scss`:

```scss
.carousel-button {
  background: rgba(0, 0, 0, 0.6);  // Cambia este color
  border: 2px solid rgba(255, 255, 255, 0.3);
}
```

### **Ajustar tamaño del modal**

```scss
.modal-container {
  max-width: 900px;  // Cambia este valor (ej: 1200px)
}
```

### **Personalizar indicadores**

```scss
.carousel-indicator {
  background: rgba(255, 255, 255, 0.4);
  
  &.active {
    background: white;  // Color del indicador activo
  }
}
```

---

## 🧪 Cómo Probar Durante Desarrollo

### **Opción 1: Resetear desde la consola**

```javascript
// En la consola del navegador
sessionStorage.clear();
localStorage.clear();
location.reload();
```

### **Opción 2: Modo incógnito**

Abre la app en ventana incógnita - el modal aparecerá siempre.

### **Opción 3: Usar ng serve con hot reload**

El servidor ya está corriendo, solo guarda los cambios y se recargará automáticamente.

---

## 🏗️ Compilación para Producción

```bash
cd front
npm run build
```

El carrusel se incluirá **automáticamente** en el bundle de producción (`dist/`).

---

## 📱 Diseño Responsive

El carrusel se adapta automáticamente:

- 💻 **Desktop**: Botones grandes, modal con margen
- 📱 **Móvil**: Pantalla completa, botones más pequeños
- 🖥️ **Tablet**: Diseño intermedio adaptativo

---

## ♿ Accesibilidad

✅ **Focus Trap** - Mantiene el foco dentro del modal  
✅ **ARIA Labels** - Etiquetas descriptivas para lectores de pantalla  
✅ **Navegación por teclado** - Flechas ← →, Tab, ESC  
✅ **Indicadores semánticos** - `aria-current` en slide activo  
✅ **Alto contraste** - Adaptación automática para accesibilidad  
✅ **Reducir movimiento** - Respeta `prefers-reduced-motion`  

---

## 🔧 Comportamiento Detallado

### **Auto-avance**

1. Inicia automáticamente al mostrar el modal
2. Cambia de slide cada 3 segundos
3. Se detiene cuando el usuario hace click en la imagen
4. Se reinicia cuando el usuario navega manualmente con flechas

### **Pausa al hacer click**

```typescript
onImageClick(): void {
  this.pauseAutoAdvance();  // Detiene el timer
}
```

### **Continuar con flechas**

```typescript
nextSlide(): void {
  this.currentSlide = (this.currentSlide + 1) % this.images.length;
  this.pauseAutoAdvance();   // Pausa el auto-avance actual
  this.startAutoAdvance();   // Reinicia un nuevo timer
}
```

---

## 📊 Recomendaciones de Imágenes

| Aspecto | Recomendación |
|---------|---------------|
| **Formato** | JPG, PNG, WebP |
| **Tamaño recomendado** | 1200x600 px (ratio 2:1) |
| **Peso máximo** | 500 KB por imagen |
| **Cantidad** | 3-5 imágenes ideal |
| **Contenido** | Diseños con texto legible y alto contraste |

---

## 🛠️ Solución de Problemas

### **El carrusel no avanza automáticamente**

Verifica que `autoAdvanceInterval` esté configurado correctamente en el CONFIG.

### **Las imágenes no se muestran**

1. Verifica que las rutas sean correctas: `assets/img/nombre.jpg`
2. Asegúrate de que las imágenes estén en `front/src/assets/img/`
3. Revisa la consola del navegador para errores

### **El auto-avance no se pausa al hacer click**

Asegúrate de que el evento `(click)="onImageClick()"` esté en el `<img>` tag del HTML.

---

## 📝 Notas Técnicas

- ⚠️ El componente es **standalone**, no requiere módulos
- ⚠️ Compatible con Angular 17+
- ⚠️ Usa `setInterval` para el auto-avance (limpiado en `ngOnDestroy`)
- ⚠️ El z-index es 9999 (sobre cualquier contenido)
- ⚠️ Las imágenes usan `loading="lazy"` para optimización

---

## 🎯 Ejemplo de Uso Completo

```typescript
// 1. Agregar tus imágenes en front/src/assets/img/
//    - comunicado-1.jpg
//    - comunicado-2.jpg
//    - comunicado-3.jpg

// 2. Actualizar el array en modal-overlay.component.ts
images = [
  {
    src: 'assets/img/comunicado-1.jpg',
    alt: 'Nueva convocatoria de becas 2026',
    fallback: 'https://via.placeholder.com/800x500/1e3a8a/ffffff?text=Becas'
  },
  {
    src: 'assets/img/comunicado-2.jpg',
    alt: 'Programa de apoyo psicológico',
    fallback: 'https://via.placeholder.com/800x500/3b82f6/ffffff?text=Psicología'
  },
  {
    src: 'assets/img/comunicado-3.jpg',
    alt: 'Actividades culturales de enero',
    fallback: 'https://via.placeholder.com/800x500/60a5fa/ffffff?text=Cultura'
  }
];

// 3. Opcional: Ajustar tiempo de auto-avance
autoAdvanceInterval: 5000  // 5 segundos

// 4. Guardar y compilar
// npm run build
```

---

**Desarrollado para la Dirección de Bienestar Universitario** 🎓  
**Versión:** Carrusel con Auto-Avance 3s
