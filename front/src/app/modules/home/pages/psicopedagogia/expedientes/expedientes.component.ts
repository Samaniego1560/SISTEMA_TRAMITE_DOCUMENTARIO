import { NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { DatePipe } from '@angular/common';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ModalSrqComponent } from "../modal-srq/modal-srq.component";
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Answer } from '../../../../../core/models/psicopedagogia/Answer';
import { ManagerService } from '../../../../../core/services/manager/manager.service';

@Component({
    // selector: 'app-resultado-cuestionario',
    standalone: true,
    imports: [
        NgFor,
        NgIf,
        FormsModule,
        ReactiveFormsModule,
        ToastComponent,
        ModalSrqComponent,
        ModalSrqComponent
    ],
    templateUrl: './expedientes.component.html',
    styleUrl: './expedientes.component.scss',
    providers: [DatePipe, ToastService]
})
export class ExpedientesComponent {
    private _subscriptions: Subscription = new Subscription();
    searchForm: FormGroup;
    fullName = '';
    srqList: any = [];
    isLoading: boolean = false;
    mostrarResultados: boolean = false;

    objParticipante: any;
    objItemSrq = {}
    respuestas: Answer[] = [];

    constructor(
        private _managerService: ManagerService,
        private fb: FormBuilder,
        private datePipe: DatePipe,
        private _toastService: ToastService
    ) {
        this.searchForm = this.fb.group({
            filterType: ['', Validators.required],
            filterValue: ['', Validators.required],
            fechaInicio: ['', Validators.required],
            fechaFin: ['', Validators.required],
        });
    }

    getHistorial() {
        const { filterType, filterValue } = this.searchForm.value;
        const dni = filterType === 'dni' ? filterValue : null;
        const apellido = filterType !== 'dni' ? filterValue : null;

        const payload = {
            dni: dni,
            apellido: apellido,
            fecha_inicio: this.searchForm.value.fechaInicio,
            fecha_fin: this.searchForm.value.fechaFin
        }

        if (!filterType || !filterValue) {
            this._toastService.add({ type: 'warning', message: "El campo fintro y valor son obligatorios" });
            return;
        }

        if ((payload.fecha_inicio && !payload.fecha_fin) || (payload.fecha_fin && !payload.fecha_inicio)) {
            this._toastService.add({ type: 'warning', message: "Ingrese los dos rango de fechas" });
            return;
        }

        if (payload.fecha_inicio && payload.fecha_fin) {
            const inicio = new Date(payload.fecha_inicio);
            const fin = new Date(payload.fecha_fin);
            if (inicio > fin) {
                this._toastService.add({ type: 'warning', message: "La fecha inicio no debe ser mayor a la fecha fin" });
                return;
            }
        }
        console.log(JSON.stringify(payload))
        this.srqList = [];
        this.isLoading = true;
        this._subscriptions.add(
            this._managerService.getHistorial(payload).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    const dataList = response.detalle;
                    console.log(JSON.stringify(dataList))
                    if (dataList.length > 0) {
                        this.fullName = `${dataList[0].nombre_estudiante} ${dataList[0].apellido_estudiante}`
                        for (let i = dataList.length - 1, index = 1; i >= 0; i--, index++) {
                            const tempName = dataList[i].es_srq ? "REGISTRO S.R.Q" : "ATENCIÓN"
                            let resul = {
                                id: index,
                                name: `${tempName} - ${index}`,
                                date: this.datePipe.transform(dataList[i].fecha_respuesta, 'dd/MM/yyyy hh:mm a'),
                                id_participante: dataList[i].id_participante,
                                data: dataList[i]
                            }
                            this.srqList.push(resul);
                        }
                    } else {
                        this._toastService.add({ type: 'warning', message: "No se encontraron registros para los filtros proporcionados" });
                    }

                },
                error: (error: any) => {
                    this._toastService.add({ type: 'warning', message: "No se encontraron registros para los filtros proporcionados" });
                    this.isLoading = false;
                }
            })
        );

    }


    handleClickViewInfoSRQ(idSrq: number) {
        this.isLoading = true;
        const { filterType, filterValue } = this.searchForm.value;
        const dni = filterType === 'dni' ? filterValue : null;
        const filters = {
            dni: dni
        };
        this.respuestas = [];
        this._subscriptions.add(
            this._managerService.searchQuestionnaires(filters).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    this.objItemSrq = this.srqList[idSrq - 1];
                    this.objParticipante = response.detalle[0];
                    this._subscriptions.add(
                        this._managerService.getResultAnswers(this.objParticipante.id, idSrq).subscribe({
                            next: (response: any) => {
                                this.isLoading = false;
                                const resultData = response.detalle;
                                if(resultData){
                                    resultData.forEach((item: any) => {
                                        const res: Answer = {
                                            texto: item.texto_pregunta,
                                            respuesta: item.respuesta.toUpperCase()
                                        };
                                        this.respuestas.push(res);
                                    }); 
                                }
                                this.mostrarResultados = true;
                            },
                            error: (error: any) => {
                                console.error('Error al obtener las respuestas del estudiante:', error);
                                this.isLoading = false;
                                this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener los resultados de la encuesta' });
                            }
                        })
                    );
                },
                error: (error: any) => {
                    this._toastService.add({ type: 'warning', message: "No se puedo recuperar al estudiante" });
                    this.isLoading = false;
                }
            })
        );
    }

    resetForm() {
        this.searchForm = this.fb.group({
            filterType: ['',],
            filterValue: ['',],
            fechaInicio: ['',],
            fechaFin: ['',],
        });
    }

    changeHistoryData(data: boolean){
        this.mostrarResultados = data
    }

}
