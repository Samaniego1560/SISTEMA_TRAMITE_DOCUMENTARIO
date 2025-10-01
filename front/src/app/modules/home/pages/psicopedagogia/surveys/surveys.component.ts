import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { NgClass, NgFor, NgIf } from '@angular/common';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../../core/services/toast/toast.service';

interface Survey {
    id: number;
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    status: string;
}

@Component({
    selector: 'app-surveys',
    standalone: true,
    imports: [NgFor, ReactiveFormsModule, NgClass, NgIf, ToastComponent],
    templateUrl: './surveys.component.html',
    styleUrls: ['./surveys.component.scss'],
    providers: [ManagerService, ToastService]
})
export class SurveysComponent implements OnInit {
    surveys: Survey[] = [];
    surveyForm: FormGroup;
    editMode = false;
    currentSurveyId!: number;
    showModal = false;
    isLoading = true;

    constructor(private fb: FormBuilder, private _managerService: ManagerService, private _toastService: ToastService,) {
        this.surveyForm = this.fb.group({
            name: [{ value: 'S.R.Q', disabled: true }, Validators.required],
            description: ['', Validators.required],
            startDate: ['', Validators.required],
            endDate: ['', Validators.required],
        });
    }

    ngOnInit(): void {
        this.getSurveys();
    }

    updateSurvey() {
        const today = new Date().toISOString().split('T')[0];
        const outdatedSurveys = this.surveys.filter(survey => survey.endDate < today && survey.status.toLowerCase() === 'activa')[0];
        if (outdatedSurveys) {
            this.isLoading = true;
            const payload = {
                estado: 'inactiva'
            }
            this._managerService.updateSurvey(payload, outdatedSurveys.id).subscribe({
                next: (response: any) => {
                    const { error } = response;
                    if (!error) {
                        this.surveys = this.surveys.map((item: any) => {
                            if (outdatedSurveys.id == item.id) {
                                item.status = 'inactiva';
                            }
                            return item;
                        })
                        this.isLoading = false;
                    }
                },
                error: (error: any) => {
                    console.error('Error al obtener encuestas:', error);
                    this.isLoading = false;
                }
            });
        } else {

        }

    }

    getSurveys() {
        this._managerService.getSurveys().subscribe({
            next: (response: any) => {
                const { error, detalle } = response;
                if (!error && detalle) {
                    this.surveys = detalle.map((survey: any) => ({
                        id: survey.id_encuesta,
                        name: survey.nombre_encuesta,
                        description: survey.descripcion,
                        startDate: survey.fecha_inicio.split('T')[0],
                        endDate: survey.fecha_fin.split('T')[0],
                        status: survey.estado
                    }));
                }
                this.updateSurvey();
                this.isLoading = false;
            },
            error: (error: any) => {
                console.error('Error al obtener encuestas:', error);
                this.isLoading = false;
            }
        });
    }

    saveSurvey(data: any) {
        const payload = {
            nombre_encuesta: "S.R.Q",
            descripcion: data.description,
            estado: "activa",
            fecha_inicio: data.startDate,
            fecha_fin: data.endDate,
        }
        this.isLoading = true;
        this._managerService.createSurvey(payload).subscribe({
            next: (response: any) => {
                const { error, detalle } = response;
                if (!error) {
                    this.isLoading = false;
                    data.id = detalle
                    this.surveys.push(data);
                }
                this.updateSurvey();
            },
            error: (error: any) => {
                console.error('Error al obtener encuestas:', error);
                this.isLoading = false;
            }
        });
    }

    openModal() {
        const hasActiveSurvey = this.surveys.some(survey => survey.status.toLowerCase() === 'activa');
        if (!this.editMode && hasActiveSurvey) {
            this._toastService.add({
                type: 'warning',
                message: 'No puede crear más encuestas S.R.Q porque tiene una activa'
            });
        } else {
            this.showModal = true;
            this.surveyForm.get('startDate')?.enable();
            this.surveyForm.get('endDate')?.enable();
        }
    }

    closeModal() {
        this.showModal = false;
        this.surveyForm.reset({ name: 'S.R.Q' });
        this.editMode = false;
    }

    handleProcessAction(){
        const payload = {
            descripcion: this.surveyForm.get('description')?.value,
            fecha_fin: this.surveyForm.get('endDate')?.value
        }
        this.isLoading = true;
        this._managerService.updateSurvey(payload, this.currentSurveyId).subscribe({
            next: (response: any) => {
                this.isLoading = false;
                if(!response.error){
                    this.surveys = this.surveys.map((item: any) => {
                        if(item.id === this.currentSurveyId){
                            item.description = payload.descripcion;
                            item.endDate = payload.fecha_fin;
                        }
                        return item;
                    })
                    this.editMode = false;
                    this.currentSurveyId = 0;
                }
            },
            error: (error: any) => {
                console.error('Error al obtener encuestas:', error);
                this.isLoading = false;
            }
        });
    }

    addSurvey() {
        if (this.surveyForm.valid) {
            if (this.editMode) {
                this.handleProcessAction()
            } else {
                const newSurvey: Survey = {
                    ...this.surveyForm.value,
                    name: 'S.R.Q',
                    status: 'activa'
                };
                this.saveSurvey(newSurvey);

            }
            this.surveyForm.reset({ name: 'S.R.Q' });
            this.closeModal();
        }
    }

    editSurvey(survey: Survey) {
        const today = new Date().toISOString().split('T')[0];
        this.surveyForm.setValue({
            name: 'S.R.Q',
            description: survey.description,
            startDate: survey.startDate,
            endDate: survey.endDate,
        });
        const endDateValue = this.surveyForm.get('endDate')?.value;
        this.surveyForm.get('startDate')?.disable();
        if (endDateValue && endDateValue < today) {
            this.surveyForm.get('endDate')?.disable();
        }
        this.editMode = true;
        this.currentSurveyId = survey.id;
        this.showModal = true;
    }
}
