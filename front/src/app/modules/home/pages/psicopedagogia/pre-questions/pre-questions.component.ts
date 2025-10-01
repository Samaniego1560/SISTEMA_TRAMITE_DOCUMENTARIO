import { NgIf } from '@angular/common';
import { Component, Output, EventEmitter } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { DataInit } from '../interfaces/data-init';
;

@Component({
    selector: 'app-pre-questions',
    standalone: true,
    imports: [FormsModule, NgIf, ToastComponent],
    templateUrl: './pre-questions.component.html',
    styleUrl: './pre-questions.component.scss',
    providers: [ToastService]
})
export class PreQuestionsComponent {
    @Output() getData = new EventEmitter<DataInit>();
    numDocumento: string = '';
    hasActiveSurvey: boolean = false;
    @Output() sendIdEncuesta = new EventEmitter<number>();
    idEncuesta!: number;
    isLoading: boolean = true;
    isButtonEnabled: boolean = false;
    private _subscriptions: Subscription = new Subscription();
    typeUser: string = 'Alumno';
    fistName: string = '';
    lastName: string = '';

    public constructor(private _managerService: ManagerService, private _toastService: ToastService) {
        this.validateActiveSRQ()
    }

    getDataQuestion() {
        this.isLoading = true;
        this._subscriptions.add(
            this._managerService.hasStudentResponded(this.numDocumento, this.idEncuesta).subscribe({
                next: (response: any) => {
                    const {respondida } = response.detalle;
                    if(respondida === true){
                        this._toastService.add({type: 'warning', message: "Usted ya ha respondido la encuesta"});
                    }else{
                        const data: DataInit = {
                            numDocumento: this.numDocumento,
                            typeUser:  this.typeUser,
                            fistName: this.fistName,
                            lastName: this.lastName
                        }
                        this.getData.emit(data);
                        this.sendIdEncuesta.emit(this.idEncuesta)
                    }
                    this.isLoading = false;
                },
                error: (error: any) => {
                    console.error('Error al realizar la petición - encuestas:', error);
                    this.hasActiveSurvey = false;
                    this.isLoading = false;
                }
            })
        );

    }

    validateInput(): void {
        this.isButtonEnabled = this.numDocumento.length === 8;
        if(this.typeUser !== 'Alumno'){
            this.isButtonEnabled = this.isButtonEnabled && this.fistName.trim().length > 1 && this.lastName.trim().length > 1;
        }
    }

    onUserTypeChange(){
        this.validateInput();
    }

    validateActiveSRQ() {
        this._subscriptions.add(
            this._managerService.hasActiveSRQ().subscribe({
                next: (response: any) => {
                    this.hasActiveSurvey = response?.detalle.activa || false;
                    this.idEncuesta = response?.detalle.id_encuesta;
                    this.isLoading = false;
                },
                error: (error: any) => {
                    console.error('Error al realizar la petición - encuestas:', error);
                    this.hasActiveSurvey = false;
                    this.isLoading = false;
                }
            })
        );
    }
}
