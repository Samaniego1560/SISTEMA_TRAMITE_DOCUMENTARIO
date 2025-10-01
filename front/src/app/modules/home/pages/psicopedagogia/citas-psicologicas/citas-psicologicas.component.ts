import { NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

import { Subscription } from 'rxjs';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { ManagerService } from '../../../../../core/services/manager/manager.service';

@Component({
    selector: 'app-citas-psicologicas',
    standalone: true,
    imports: [FormsModule, ReactiveFormsModule, NgIf, ToastComponent],
    templateUrl: './citas-psicologicas.component.html',
    styleUrl: './citas-psicologicas.component.scss',
    providers: [ToastService]
})
export class CitasPsicologicasComponent {
    citaForm: FormGroup;
    errorHorario: string = '';
    private _subscriptions: Subscription = new Subscription();
    isLoading: boolean = false

    constructor(
        private fb: FormBuilder, 
        private _toastService: ToastService,
        private _managerService: ManagerService,
    ) {
        this.citaForm = this.fb.group({
            dni: ['', [Validators.required, Validators.pattern('^[0-9]{8}$')]],
            nombre: ['', Validators.required],
            apellido: ['', Validators.required],
            facultad: ['', Validators.required],
            hora_inicio: ['', Validators.required],
        });

        this.citaForm.get('hora_inicio')?.valueChanges.subscribe((value) => {
            this.validarHorario(value);
        });
    }

    validarHorario(fecha: string) {
        const date = new Date(fecha);
        const day = date.getDay();
        const hour = date.getHours();
        const minute = date.getMinutes();

        if (day === 0 || day === 6) {
            this.errorHorario = 'No se atiende los sábados y domingos';
        } else if (hour < 8 || hour > 13 || (hour === 13 && minute > 0)) {
            this.errorHorario = 'El horario de atención es de 8:00 AM a 1:00 PM';
        } else if(date <= new Date()){
            this.errorHorario = 'Selecciona un horario a partir de mañana';
        }else{
            this.errorHorario = '';
        }
    }

    submitForm() {
        if (this.citaForm.valid && this.errorHorario.length === 0) {
            console.log('Formulario enviado:', this.citaForm.value);
            const payload = {...this.citaForm.value};
            payload.hora_inicio = `${payload.hora_inicio}:00Z`
            this.isLoading = true;
            this.citaForm.get('')
            this._subscriptions.add(
                this._managerService.createCita(payload).subscribe({
                    next: (response: any) => {
                        this.isLoading = false;
                        this._toastService.add({ type: 'success', message: 'La cita se registro correctamente' });
                        console.log(response)
                    },
                    error: (error: any) => {
                        console.log(error)
                        this.isLoading = false;
                        this._toastService.add({ type: 'error', message: error.error.msg });
                        console.error('Error al crear la cita:', error);
                    }
                })
            )
        }else if(this.errorHorario.length > 0){
            this._toastService.add({type: 'warning', message: this.errorHorario});
        }else{
            this._toastService.add({type: 'warning', message: "Por favor, completa todos los campos correctamente"});
        }
    }
}
