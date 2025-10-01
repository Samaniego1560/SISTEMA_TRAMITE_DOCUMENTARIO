import { NgFor, NgIf } from '@angular/common';
import { Component, EventEmitter, inject, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { ManagerService } from '../../../../../core/services/manager/manager.service';

@Component({
    selector: 'app-registro-post-atencion',
    standalone: true,
    imports: [NgIf, FormsModule, NgFor, ToastComponent],
    templateUrl: './registro-post-atencion.component.html',
    styleUrl: './registro-post-atencion.component.scss',
    providers: [ToastService]
})
export class RegistroPostAtencionComponent {

    hasError: boolean = false;
    formData = {
        motivoConsulta: '',
        situacionActual: '',
        diagnosticoPresuntivo: '',
        otrosProcedimientos: '',
        notasExterno: '',
        instrumentos: '',
        resultados: ''
    };
    diagnosticos: any[] = [];

    mostrarModal = false;
    codigoDiagnostico = '';
    nombreDiagnostico = '';
    sendSuccess = false;
    isLoading = false;
    @Output() sendSuccessChange = new EventEmitter<boolean>();
    @Output() sendData = new EventEmitter<any>();
    private _subscriptions: Subscription = new Subscription();

    constructor(
        private _managerService: ManagerService,
        private _toastService: ToastService,
    ) {
        this.getAllDiagnostico()
    }

    onSubmit(): void {
        this.hasError = !this.formData.motivoConsulta || !this.formData.situacionActual ||
            !this.formData.diagnosticoPresuntivo || !this.formData.otrosProcedimientos;

        if (!this.hasError) {
            console.log('Form Submitted', this.formData);
            this.sendSuccess = true;
            this.sendSuccessChange.emit(this.sendSuccess);
            this.sendData.emit(this.formData)
        }
    }

    agregarDiagnostico() {
        this.mostrarModal = true;
    }

    guardarDiagnostico() {
        if (this.codigoDiagnostico.trim() == '' || this.nombreDiagnostico.trim() == '') {
            this._toastService.add({
                type: 'warning',
                message: 'Completa los dos campos!',
                life: 5000
            })
        } else {
            this.isLoading = true;
            const payload = {
                codigo : this.codigoDiagnostico,
                nombre: this.nombreDiagnostico
            }
            let diagnosticoId = 0;
            this.cerrarModal();
            this._subscriptions.add(
                this._managerService.createDiagnostico(payload).subscribe({
                    next: (response: any) => {
                        this.isLoading = false;
                        if (response && response.detalle) {
                            diagnosticoId = response.detalle;
                        }
                        this.diagnosticos.push({
                            id: diagnosticoId,
                            nombre: this.nombreDiagnostico
                        })
                        this.codigoDiagnostico = '';
                        this.nombreDiagnostico = ''
                        this.isLoading = false;
                    },
                    error: (error: any) => {
                        this.isLoading = false;
                        console.error('Error al obtener las respuestas del estudiante:', error);
                    }
                })
            );
        }
    }

    cerrarModal() {
        this.mostrarModal = false;
    }

    getAllDiagnostico() {
        this.isLoading = true;
        this._subscriptions.add(
            this._managerService.getAllDiagnostico().subscribe({
                next: (response: any) => {
                    if (response && response.detalle) {
                        const detalle = response.detalle;
                        this.diagnosticos = detalle
                    }
                    this.isLoading = false;
                },
                error: (error: any) => {
                    this.isLoading = false;
                    console.error('Error al obtener las respuestas del estudiante:', error);
                }
            })
        );
    }
}
