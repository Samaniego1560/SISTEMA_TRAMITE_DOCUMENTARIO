import { Component, effect } from '@angular/core';
import { BlockUiComponent } from "../../../../core/ui/block-ui/block-ui.component";
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from "@angular/forms";
import { DatePipe, NgForOf, NgIf } from "@angular/common";
import { ToastComponent } from "../../../../core/ui/toast/toast.component";
import { ToastService } from "../../../../core/services/toast/toast.service";
import { CarouselSettingService } from "../../../../core/services/carousel-setting/carousel-setting.service";
import { take } from "rxjs";
import { IResponse } from "../../../../core/models/response";
import { ICarouselSetting } from "../../../../core/models/carousel-setting";

@Component({
    selector: 'app-carousel-setting',
    standalone: true,
    imports: [
        BlockUiComponent,
        FormsModule,
        NgForOf,
        NgIf,
        ReactiveFormsModule,
        ToastComponent,
        DatePipe
    ],
    templateUrl: './carousel-setting.component.html',
    styleUrl: './carousel-setting.component.scss',
    providers: [ToastService]
})
export class CarouselSettingComponent {
    protected formCarousel: FormGroup;
    public isLoading: boolean = false;
    public carouselSettings: ICarouselSetting[] = [];
    public selectedFile: File | null = null;
    public imagePreview: string | null = null;
    public editingId: number | null = null;

    constructor(private _fb: FormBuilder,
        private _toastService: ToastService,
        private _carouselSettingService: CarouselSettingService) {
        this.formCarousel = this._fb.group({
            image: [null, []],
            title: ['', []],
            description: ['', []],
            button_text: ['', []],
            button_link: ['', []],
            is_enabled: [true, []],
            order: [0, [Validators.min(0)]]
        });

        effect(() => {
            this.loadCarouselSettings().then();
        });
    }

    private async loadCarouselSettings() {
        this.isLoading = true;
        const response = await this.getCarouselSettingsPromise();
        this.isLoading = false;

        if (response.error) {
            this.eventMessage('error', response.msg);
            return;
        }

        if (response.data) {
            this.carouselSettings = response.data;
        }
    }

    public onFileSelected(event: any) {
        const file = event.target.files[0];
        if (file) {
            this.selectedFile = file;

            // Create image preview
            const reader = new FileReader();
            reader.onload = (e: any) => {
                this.imagePreview = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    }

    public async saveForm() {
        if (!this.selectedFile && !this.editingId) {
            this.eventMessage('info', 'Por favor seleccione una imagen');
            return;
        }

        const formData = new FormData();

        if (this.selectedFile) {
            formData.append('image', this.selectedFile);
        }

        const formValue = this.formCarousel.value;
        if (formValue.title) formData.append('title', formValue.title);
        if (formValue.description) formData.append('description', formValue.description);
        if (formValue.button_text) formData.append('button_text', formValue.button_text);
        if (formValue.button_link) formData.append('button_link', formValue.button_link);
        formData.append('is_enabled', formValue.is_enabled ? '1' : '0');
        formData.append('order', formValue.order.toString());

        this.isLoading = true;

        let response;
        if (this.editingId) {
            response = await this.updateCarouselSettingPromise(this.editingId, formData);
        } else {
            response = await this.createCarouselSettingPromise(formData);
        }

        this.isLoading = false;

        if (response.error) {
            this.eventMessage('error', response.msg);
            return;
        }

        this.eventMessage('success', this.editingId ? 'Carrusel actualizado correctamente' : 'Carrusel creado correctamente');
        this.resetForm();
        await this.loadCarouselSettings();
    }

    public editCarousel(carousel: ICarouselSetting) {
        this.editingId = carousel.id || null;
        this.formCarousel.patchValue({
            title: carousel.title || '',
            description: carousel.description || '',
            button_text: carousel.button_text || '',
            button_link: carousel.button_link || '',
            is_enabled: carousel.is_enabled,
            order: carousel.order
        });
        this.imagePreview = carousel.image_url || null;

        // Scroll to form
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    public async deleteCarousel(id: number) {
        if (!confirm('¿Está seguro de eliminar este carrusel?')) {
            return;
        }

        this.isLoading = true;
        const response = await this.deleteCarouselSettingPromise(id);
        this.isLoading = false;

        if (response.error) {
            this.eventMessage('error', response.msg);
            return;
        }

        this.eventMessage('success', 'Carrusel eliminado correctamente');
        await this.loadCarouselSettings();
    }

    public async toggleEnabled(carousel: ICarouselSetting) {
        const formData = new FormData();
        formData.append('is_enabled', carousel.is_enabled ? '0' : '1');
        formData.append('order', carousel.order.toString());

        this.isLoading = true;
        const response = await this.updateCarouselSettingPromise(carousel.id!, formData);
        this.isLoading = false;

        if (response.error) {
            this.eventMessage('error', response.msg);
            return;
        }

        await this.loadCarouselSettings();
    }

    public resetForm() {
        this.formCarousel.reset({
            title: '',
            description: '',
            button_text: '',
            button_link: '',
            is_enabled: true,
            order: this.carouselSettings.length
        });
        this.selectedFile = null;
        this.imagePreview = null;
        this.editingId = null;
    }

    private createCarouselSettingPromise(formData: FormData): Promise<{ data: ICarouselSetting | null, error: boolean, msg: string }> {
        return new Promise((resolve) => {
            this._carouselSettingService.createCarouselSetting(formData).pipe(take(1)).subscribe({
                next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
                error: (err) => resolve({ data: null, error: true, msg: err.error?.mensaje || 'Error al crear carrusel.' })
            });
        });
    }

    private updateCarouselSettingPromise(id: number, formData: FormData): Promise<{ data: ICarouselSetting | null, error: boolean, msg: string }> {
        return new Promise((resolve) => {
            this._carouselSettingService.updateCarouselSetting(id, formData).pipe(take(1)).subscribe({
                next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
                error: (err) => resolve({ data: null, error: true, msg: err.error?.mensaje || 'Error al actualizar carrusel.' })
            });
        });
    }

    private getCarouselSettingsPromise(): Promise<{ data: ICarouselSetting[] | null, error: boolean, msg: string }> {
        return new Promise((resolve) => {
            this._carouselSettingService.getCarouselSettings().pipe(take(1)).subscribe({
                next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
                error: () => resolve({ data: null, error: true, msg: 'Error al obtener configuración de carrusel.' })
            });
        });
    }

    private deleteCarouselSettingPromise(id: number): Promise<{ data: null, error: boolean, msg: string }> {
        return new Promise((resolve) => {
            this._carouselSettingService.deleteCarouselSetting(id).pipe(take(1)).subscribe({
                next: (res: IResponse<any>) => resolve({ data: null, error: false, msg: '' }),
                error: () => resolve({ data: null, error: true, msg: 'Error al eliminar carrusel.' })
            });
        });
    }

    public eventMessage(type: "error" | "success" | "warning" | "info", message: string): void {
        this._toastService.add({
            type: type,
            message: message,
        });
    }
}
