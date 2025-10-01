import { NgClass, NgFor, NgForOf, NgIf } from '@angular/common';
import { Component, Output, EventEmitter, Input } from '@angular/core';
import { Subscription } from 'rxjs';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { PaginationComponent } from '../../pagination/pagination';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Student } from '../../../../../core/models/psicopedagogia/students';
import { Answer } from '../../../../../core/models/psicopedagogia/Answer';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { generateSRQPDF, getBase64ImageFromUrl } from '../utils/pdf-sqr-util';


@Component({
    selector: 'app-resultado-cuestionario',
    standalone: true,
    imports: [
        NgClass,
        NgFor,
        NgForOf,
        NgIf,
        FormsModule,
        ToastComponent,
        ReactiveFormsModule,
        PaginationComponent],
    templateUrl: './resultado-cuestionario.component.html',
    styleUrl: './resultado-cuestionario.component.scss',
    providers: [ToastService]
})
export class ResultadoCuestionarioComponent {
    students: Student[] = [];
    infoStudent!: Student;
    respuestas: Answer[] = [];
    @Output() close = new EventEmitter<void>();
    // showResultSRQ: boolean = false;
    @Input() showResultSRQ: boolean = false;
    private _subscriptions: Subscription = new Subscription();
    isLoading: boolean = true;
    totalSi: number = 0;
    totalNo: number = 0;
    diagnostico: string = '';
    diagnosticoClass: string = '';
    studentMessage: string = '';
    selectedAction: string | null = null;
    showFullDetails: boolean = false;

    searchForm: FormGroup;
    currentPage = 1;
    totalPages = 0;
    responseData: any;

    constructor(
        private _managerService: ManagerService,
        private _toastService: ToastService,
        private fb: FormBuilder
    ) {
        this.getResultQuestionnaire();
        this.searchForm = this.fb.group({
            filterType: ['', Validators.required],
            filterValue: ['', Validators.required]
        });
    }

    onPageChange(page: number): void {
        this.currentPage = page;
        console.log('Página cambiada a:', page);
        this.isLoading = true;
        this.students = [];
        this.getResultQuestionnaire();
    }

    getButtonText(): string {
        if (this.studentMessage.trim() && this.selectedAction) {
            return 'Guardar y Enviar';
        } else if (this.selectedAction) {
            return 'Guardar';
        }
        return 'Seleccione una opción';
    }

    calcularResultados() {
        // Contar respuestas "Sí" y "No"
        this.totalSi = this.respuestas.filter(p => p.respuesta === 'SI').length;
        this.totalNo = this.respuestas.filter(p => p.respuesta === 'NO').length;

        // Determinar diagnóstico
        if (this.totalSi >= 7) {
            this.diagnostico = 'Grave';
            this.diagnosticoClass = 'text-red-500';
        } else if (this.totalSi >= 5 && this.totalSi <= 6) {
            this.diagnostico = 'Moderado';
            this.diagnosticoClass = 'text-yellow-500';
        } else {
            this.diagnostico = 'Estable';
            this.diagnosticoClass = 'text-green-500';
        }
    }

    public getResultQuestionnaire(): void {
        this._subscriptions.add(
            this._managerService.getResultQuestionnaire(this.currentPage).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    const resultData = response.detalle.items;
                    this.totalPages = response.detalle.totalPages;
                    resultData.forEach((item: any) => {
                        const fullName = `${item.nombre} ${item.apellido}`;
                        const res: Student = {
                            id: item.id,
                            numeroAtencion: item.numero_atencion,
                            studentId: item.student_id,
                            tipoPaciente: item.tipo_participante,
                            code: item.codigo_estudiante,
                            dni: item.dni,
                            numPhone: item.num_telefono,
                            school: item.escuela,
                            diagnostico: item.diagnostico || 'No Aplica',
                            status: item.estado_evaluacion || 'No Aplica',
                            fullName: fullName,
                            birthDate: item.fecha_nacimiento?.split("T")[0],
                            age: item.edad,
                            previousSchool: item.colegio_procedencia,
                            birthPlace: item.lugar_nacimiento,
                            entryMode: item.modalidad_ingreso,
                            currentLivingSituation: item.con_quienes_vive_actualmente,
                            universityEntryYear: item.anio_ingreso,
                            currentSemester: item.semestre_cursa,
                            evaluationDate: item.created_at.split("T")[0]
                        };

                        this.students.push(res);
                    });
                },
                error: (error: any) => {
                    console.error('Error al obtener los resultados:', error);
                    this.isLoading = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener los resultados de la encuesta' });
                }
            })
        );
    }

    marcarComoRevisado() {
        alert('El cuestionario ha sido marcado como revisado.');
    }

    handleViewResultSRQ(studentId: number) {
        this.isLoading = true;
        this.showFullDetails = false;
        this.infoStudent = this.students.find(student => student.id === studentId)!;
        if (this.infoStudent.status !== 'Aprobado') {
            this.showFullDetails = true;
        }
        this.respuestas = [];
        this._subscriptions.add(
            this._managerService.getResultAnswers(studentId, this.infoStudent.numeroAtencion).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    const resultData = response.detalle;
                    resultData.forEach((item: any) => {
                        const res: Answer = {
                            texto: item.texto_pregunta,
                            respuesta: item.respuesta.toUpperCase()
                        };
                        this.respuestas.push(res);
                    });
                    this.calcularResultados();
                },
                error: (error: any) => {
                    console.error('Error al obtener las respuestas del estudiante:', error);
                    this.isLoading = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener los resultados de la encuesta' });
                }
            })
        );
        if (this.infoStudent.diagnostico == 'Estable' && this.infoStudent.status == 'Aprobado') {
            this.showFullDetails = false;
        }
        this.showResultSRQ = true;
    }

    closeModal() {
        this.studentMessage = '';
        this.selectedAction = null;
        this.showResultSRQ = false;
    }

    handleAction() {
        this.isLoading = true;
        const dataUpdate = {
            statusEvaluation: this.selectedAction === 'observed' ? 'Observado' : 'Aprobado',
            notes: this.studentMessage,
        }
        this._subscriptions.add(
            this._managerService.updateStatusEvaluationParticipante(dataUpdate, Number(this.infoStudent.id)).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    console.log(JSON.stringify(response));
                    this._toastService.add({ type: 'success', message: 'Datos actualizados correctamente' });
                    this.students = [];
                    if (this.selectedAction === 'approved') {
                        this.showFullDetails = false;
                    }
                    this.getResultQuestionnaire();
                },
                error: (error: any) => {
                    console.error('Error al actualizar el estado:', error);
                    this.isLoading = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al actualizar los datos' });
                }
            })
        );
    }


    onSearch(): void {
        const type = this.searchForm.value.filterType;
        const value = this.searchForm.value.filterValue;
        console.log(type)
        if (type == "" || (value == "" && type != "todos")) {
            this._toastService.add({ type: 'warning', message: 'Por favor selecciona un tipo de filtro e ingrese un valor. ' });
            return;
        }
        const filters = {
            [type]: value
        };
        this.isLoading = true;
        console.log(JSON.stringify(filters))
        this._subscriptions.add(
            this._managerService.searchQuestionnaires(filters).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    console.log(JSON.stringify(response))
                    this.students = [];
                    if (response.detalle) {
                        const resultData = response.detalle;
                        resultData.forEach((item: any) => {
                            const fullName = `${item.nombre} ${item.apellido}`;
                            const res: Student = {
                                id: item.id,//tengo
                                numeroAtencion: item.numero_atencion,
                                studentId: item.student_id, //tengo
                                tipoPaciente: item.tipo_participante,
                                code: item.codigo_estudiante,//tengo
                                dni: item.dni,
                                numPhone: item.num_telefono,//tengo
                                school: item.escuela,//tengo
                                diagnostico: item.diagnostico,//tengo
                                status: item.estado_evaluacion,//tengo
                                fullName: fullName,//tengo
                                birthDate: item.fecha_nacimiento?.split("T")[0],
                                age: item.edad,
                                previousSchool: item.colegio_procedencia,//tengo
                                birthPlace: item.lugar_nacimiento,
                                entryMode: item.modalidad_ingreso,
                                currentLivingSituation: item.con_quienes_vive_actualmente,//tengo
                                universityEntryYear: item.anio_ingreso,//tengo
                                currentSemester: item.semestre_cursa,//tengo
                                evaluationDate: item.created_at.split("T")[0]//tengo

                            };

                            this.students.push(res);
                        });
                    }

                },
                error: (error: any) => {
                    console.error('Error al obtener las respuestas del estudiante:', error);
                    this.isLoading = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener los resultados de la encuesta' });
                }
            })
        );
    }

    handleDowloaderSrq(studentId: number) {
        this.isLoading = true;
        this._subscriptions.add(
            this._managerService.getPdfDowloader(studentId).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    this.generarPDF(response.detalle);
                },
                error: (error: any) => {
                    console.error('Error al obtener las respuestas del estudiante:', error);
                    this.isLoading = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al descargar el documento S.R.Q' });
                }
            })
        );
    }

    async generarPDF(data: any){
        const logoUnas = await getBase64ImageFromUrl('assets/logos/logo-psico.png');
        const logoPsico = await getBase64ImageFromUrl('assets/logos/logo-unas.png');
        generateSRQPDF(data, this._toastService, logoPsico, logoUnas);
    }
}


