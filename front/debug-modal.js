/* 
 * SCRIPT DE DEBUGGING PARA EL MODAL CAROUSEL
 * 
 * INSTRUCCIONES:
 * 1. Abre tu aplicación en el navegador (http://localhost:4200)
 * 2. Abre la Consola del Navegador (F12 > Console)
 * 3. Copia y pega ESTE SCRIPT COMPLETO en la consola
 * 4. Presiona Enter
 * 5. Lee los resultados en la consola
 */

console.log('🔍 INICIANDO DEBUG DEL MODAL CAROUSEL...\n');

// 1. Verificar Storage
console.log('📦 1. Verificando Storage:');
const sessionValue = sessionStorage.getItem('dbu_modal_last_shown');
const localValue = localStorage.getItem('dbu_modal_last_shown');
console.log('   - sessionStorage:', sessionValue || 'VACÍO');
console.log('   - localStorage:', localValue || 'VACÍO');

if (sessionValue || localValue) {
    console.log('   ⚠️ ENCONTRADO! El modal ya se mostró antes.');
    console.log('   💡 Limpiando storage...');
    sessionStorage.removeItem('dbu_modal_last_shown');
    localStorage.removeItem('dbu_modal_last_shown');
    console.log('   ✅ Storage limpiado. Recargando página...');
    setTimeout(() => location.reload(), 1000);
} else {
    console.log('   ✅ Storage limpio');
}

// 2. Verificar si el modal está en el DOM
console.log('\n📋 2. Verificando DOM:');
const modalOverlay = document.querySelector('.modal-overlay');
const modalContainer = document.querySelector('.modal-container');
const carouselSlides = document.querySelector('.carousel-slides');

console.log('   - modal-overlay:', modalOverlay ? '✅ Encontrado' : '❌ No encontrado');
console.log('   - modal-container:', modalContainer ? '✅ Encontrado' : '❌ No encontrado');
console.log('   - carousel-slides:', carouselSlides ? '✅ Encontrado' : '❌ No encontrado');

// 3. Verificar imágenes
console.log('\n🖼️ 3. Verificando imágenes:');
const slides = document.querySelectorAll('.carousel-slide');
const images = document.querySelectorAll('.carousel-image');

console.log('   - Número de slides:', slides.length);
console.log('   - Número de imágenes:', images.length);

if (images.length > 0) {
    console.log('   - Rutas de imágenes:');
    images.forEach((img, i) => {
        console.log(`     ${i + 1}. ${img.src}`);
        console.log(`        - Cargada: ${img.complete ? '✅' : '⏳ Cargando...'}`);
        console.log(`        - Error: ${img.naturalWidth === 0 ? '❌ SÍ' : '✅ NO'}`);
    });
}

// 4. Verificar clases activas
console.log('\n🎯 4. Verificando estado:');
if (modalOverlay) {
    console.log('   - Overlay visible:', modalOverlay.classList.contains('active') ? '✅ SÍ' : '❌ NO');
    console.log('   - Overlay closing:', modalOverlay.classList.contains('closing') ? '⚠️ SÍ' : '✅ NO');
}

if (slides.length > 0) {
    slides.forEach((slide, i) => {
        if (slide.classList.contains('active')) {
            console.log(`   - Slide activo: #${i + 1} ✅`);
        }
    });
}

// 5. Verificar CSS
console.log('\n🎨 5. Verificando estilos:');
if (modalOverlay) {
    const overlayStyle = window.getComputedStyle(modalOverlay);
    console.log('   - opacity:', overlayStyle.opacity);
    console.log('   - visibility:', overlayStyle.visibility);
    console.log('   - display:', overlayStyle.display);
}

if (carouselSlides) {
    const slidesStyle = window.getComputedStyle(carouselSlides);
    console.log('   - slides height:', slidesStyle.height);
    console.log('   - slides width:', slidesStyle.width);
}

// 6. Función auxiliar para forzar mostrar el modal
console.log('\n🔧 6. Funciones de ayuda disponibles:');
console.log('   - Para forzar mostrar el modal, ejecuta: mostrarModal()');
console.log('   - Para resetear y recargar: resetearModal()');

window.mostrarModal = function () {
    console.log('📢 Intentando forzar mostrar modal...');
    sessionStorage.removeItem('dbu_modal_last_shown');
    localStorage.removeItem('dbu_modal_last_shown');

    const overlay = document.querySelector('.modal-overlay');
    if (overlay) {
        overlay.classList.add('active');
        overlay.style.opacity = '1';
        overlay.style.visibility = 'visible';
        document.body.classList.add('modal-open');
        console.log('✅ Modal forzado a mostrarse');
    } else {
        console.log('❌ No se encontró el elemento .modal-overlay');
    }
};

window.resetearModal = function () {
    console.log('🔄 Reseteando modal...');
    sessionStorage.clear();
    localStorage.clear();
    setTimeout(() => location.reload(), 500);
};

console.log('\n✅ DEBUG COMPLETADO\n');
console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
console.log('💡 SOLUCIÓN RÁPIDA:');
console.log('   Si el modal no aparece, ejecuta:');
console.log('   resetearModal()');
console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');
