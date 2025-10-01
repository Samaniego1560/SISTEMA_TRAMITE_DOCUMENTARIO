import { CommonModule, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule, NgModel } from '@angular/forms';
import { Subscription } from 'rxjs';
import { v4 as uuidv4 } from 'uuid';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { EnvServiceFactory } from '../../../../../core/services/env/env.service.provider';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { RegistroAtencionesComponent } from '../registro-atenciones/registro-atenciones.component';
import { RegistroPostAtencionComponent } from '../registro-post-atencion/registro-post-atencion.component';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';

@Component({
    selector: 'app-atenciones',
    standalone: true,
    imports: [NgIf, CommonModule, RegistroAtencionesComponent, RegistroPostAtencionComponent, FormsModule, ToastComponent],
    templateUrl: './atenciones.component.html',
    styleUrl: './atenciones.component.scss',
    providers: [ToastService]
})
export class AtencionesComponent {
    currentStep: number = 1;
    selectedType: string = '';
    linkRecibido: string = '';
    isPostRegEnabled: boolean = false;
    isGenerateLinkEnabled: boolean = false;
    showModal: boolean = false;
    data: any;
    additionalData: any;
    disableRegisterAtencion = false;
    isLoading = false;
    isSuccess = false;
    private _subscriptions: Subscription = new Subscription();
    private urlFrond = EnvServiceFactory().URL_FROND;

    constructor(
        private _managerService: ManagerService,
        private _toastService: ToastService,
    ) { }

    selectUserType(type: string) {
        this.selectedType = type;
    }

    nextStep() {
        if (this.currentStep < 3) {
            this.currentStep++;
        }

    }

    previousStep() {
        if (this.currentStep > 1) {
            this.currentStep--;
        }
    }

    confirmProcess() {
        this.showModal = true;
    }

    confirmExit() {
        this.currentStep = 1;
        this.data = null;
        this.isPostRegEnabled = false;
        this.isGenerateLinkEnabled = false;
        this.linkRecibido = '';
        this.selectedType = '';
        this.showModal = false;
        this.disableRegisterAtencion = false;
    }

    close() {
        this.showModal = false;
    }

    handleProcessAdditionalData(e: any) {
        const participante = {
            tipo_participante: this.data.tipoParticipante,
            nombre: this.data.firstName,
            apellido: this.data.lastName,
            dni: this.data.dni,
            num_telefono: this.data.phone,
            con_quienes_vive_actualmente: this.data.currentLiving,
            fecha_nacimiento: this.data.birthDate,
            edad: Number( this.data.age),
            lugar_nacimiento: this.data.birthPlace,
            colegio_procedencia: this.data.school,
            anio_ingreso: Number(this.data.admissionYear),
            semestre_cursa: this.data.currentSemester,
            modalidad_ingreso: this.data.admissionMode,
            sexo: this.data.gender,
            diagnostico_id: Number(e.diagnosticoPresuntivo),
            es_srq: false,
            direccion: this.data.address,
            quien_financia_carrera: this.data.careerFinancing,
            motivo_consulta: e.motivoConsulta,
            situacion_actual: e.situacionActual,
            otros_procedimientos: e.otrosProcedimientos,
            profesion: this.data.profession,
            estado_civil: this.data.maritalStatus,
            labora_en_unas: this.data.workingAtUNAS == "Sí" ? true : false,
            grado_instruccion: this.data.degree
        }
        const payload = {
            participante: participante,
        }
        console.log(JSON.stringify(payload))
        this.isLoading = true;
        this._subscriptions.add(
            this._managerService.saveHistoryNoSrq(payload).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    this.isSuccess = true;
                },
                error: (error: any) => {
                    console.error('error guardando los registros:', error);
                    this.isSuccess = false;
                    this.isLoading = false;
                    this.disableRegisterAtencion = false;
                }
            })
        );
    }

    handleGenerateLink() {
        this.isGenerateLinkEnabled = true;
        this.data.key_url = uuidv4();
        const strData = JSON.stringify(this.data);
        const datosLimpios = strData.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
        const datos = JSON.stringify(datosLimpios);
        const base64Datos = btoa(datos);
        this.linkRecibido = `${this.urlFrond}/home/cuestionario?dataId=${base64Datos}`;
    }

    updateStep(step: number) {
        this.currentStep = step;
        console.log("Paso actualizado desde el hijo:", step);
    }

    toBase64(str: string): string {
        return btoa(new TextEncoder().encode(str).reduce((acc, byte) => acc + String.fromCharCode(byte), ''));
    }

    recibirStrData(data: string) {
        this.data = JSON.parse(data);
    }

    handleSendSuccessChange(event: boolean) {
        this.disableRegisterAtencion = event;
    }

    copyToClipboard() {
        navigator.clipboard.writeText(this.linkRecibido).then(() => {
            this._toastService.add({
                type: 'success',
                message: 'Enlace copiado al portapapeles',
                life: 5000
            })
        }).catch((err: any) => {
            console.error('Error al copiar:', err);
        });
    }
}
