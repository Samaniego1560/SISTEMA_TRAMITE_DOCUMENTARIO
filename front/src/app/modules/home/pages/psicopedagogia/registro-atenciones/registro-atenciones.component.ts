import { NgFor, NgIf } from '@angular/common';
import { Subscription } from 'rxjs';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ManagerService } from '../../../../../core/services/manager/manager.service';

@Component({
    selector: 'app-registro-atenciones',
    standalone: true,
    imports: [FormsModule, NgIf, NgFor],
    templateUrl: './registro-atenciones.component.html',
    styleUrl: './registro-atenciones.component.scss'
})
export class RegistroAtencionesComponent implements OnInit {
    @Input() typePatient!: string;
    isLoading = false;
    formData: { [key: string]: any } = {
        dni: '',
        lastName: '',
        firstName: '',
        gender: '',
        age: '',
        address: '',
        birthPlace: '',
        currentLiving: '',
        currentSemester: '',
        birthDate: '',
        school: '',
        admissionMode: '',
        admissionYear: '',
        phone: '',
        tipoParticipante: '',
        motivoConsulta: '',
        situacionActual: '',
        careerFinancing: '',
        workingAtUNAS: '',
        degree: '',
        profession: '',
        maritalStatus: '',

        consultationReason: '',
        currentSituation: '',
        presumptiveDiagnosis: '',
        otherProcedures: ''
    };
    private _subscriptions: Subscription = new Subscription();
    isConsulta: boolean = false;
    currentStep: number = 2;
    @Output() stepChanged = new EventEmitter<number>();

    ciclos: string[] = [
        "1°", "2°", "3°", "4°", "5°",
        "6°", "7°", "8°", "9°", "10°"
    ];
    diagnosticos: string[] = [
        "A", "B"
    ]
    hasError: boolean = false;
    @Output() srtData = new EventEmitter<string>();

    public constructor(private _managerService: ManagerService) { }

    ngOnInit(): void {
        const storedFormData = localStorage.getItem('formData');
        if (storedFormData) {
            this.formData = JSON.parse(storedFormData);
        } else {
            this.formData = {}; 
        }
    }

    updateField(field: string, event: Event): void {
        const inputElement = event.target as HTMLInputElement;
        const value = inputElement.value;
        this.formData[field] = value;
        localStorage.setItem('formData', JSON.stringify(this.formData));
    }

    clearFields(fieldsToIgnore: any): void {
        Object.keys(this.formData).forEach(field => {
            if (!fieldsToIgnore.includes(field)) {
                this.formData[field] = '';
            }
        });
    }

    buscarDni() {
        const dniValue = this.formData['dni'];
        this.isLoading = true;
        this.clearFields(["dni"]);
        this._subscriptions.add(
            this._managerService.getParticipantByDocument(dniValue, this.typePatient).subscribe({
                next: (response: any) => {
                    this.isLoading = false;
                    if (response && response.detalle) {
                        const detalle = response.detalle;
                        // Asignar valores de la respuesta a formData
                        this.formData = {
                            dni: detalle.dni || '',
                            firstName: detalle.nombre || detalle.nombres,
                            lastName: detalle.apellido || `${detalle.apellido_paterno} ${detalle.apellido_materno}`,
                            age: detalle.edad || '',
                            birthPlace: detalle.lugar_nacimiento || '',
                            currentLiving: detalle.con_quienes_vive_actualmente || '',
                            birthDate: this.formatDate(detalle.fecha_nacimiento),
                            school: detalle.colegio_procedencia || '',
                            admissionMode: detalle.modalidad_ingreso || '',
                            admissionYear: detalle.anio_ingreso || '',
                            phone: detalle.num_telefono || detalle.celular_estudiante,
                            tipoParticipante: this.typePatient,
                            gender: detalle.sexo,
                            address: detalle.direccion,
                            careerFinancing: detalle.quien_financia_carrera
                        };
                        localStorage.setItem('formData', JSON.stringify(this.formData));
                    }
                },
                error: (error: any) => {
                    this.clearFields('');
                    this.isLoading = false;
                    console.error('Error al obtener las respuestas del estudiante:', error);
                }
            })
        );
    }

    formatDate(dateString: string): string {
        if (!dateString) return '';
        const ddMmYyyyRegex = /^\d{2}-\d{2}-\d{4}$/;
        if (ddMmYyyyRegex.test(dateString)) {
            const parts = dateString.split('-');
            return `${parts[2]}-${parts[1]}-${parts[0]}`;;
        }
        const datePart = dateString.split('T')[0];
        const yyyymmddRegex = /^\d{4}-\d{2}-\d{2}$/;
        if (yyyymmddRegex.test(datePart)) {
            const parts = datePart.split('-');
            return `${parts[0]}-${parts[1]}-${parts[2]}`;
        }

        return '';
    }

    enviarData() {
        this.formData['tipoParticipante'] = this.typePatient;
        const datos = JSON.stringify(this.formData);
        this.srtData.emit(datos);
    }


    toggleConsulta(event: Event) {
        if (event instanceof Event) {
            const target = event.target;
            if (target instanceof HTMLInputElement) {
                this.isConsulta = target.checked;
            }
        }
    }

    nextStep() {
        this.validateFields();
        if (!this.hasError) {
            if (this.currentStep < 3) {
                this.currentStep++;
                this.stepChanged.emit(this.currentStep);
                this.enviarData();
            }
        }


    }

    validateFields() {
        const fieldsByType: { [key: string]: string[] } = {
            "Alumno": [
                'dni', 'firstName', 'lastName', 'address', 'gender', 'phone', 'age', 'birthPlace',
                'currentLiving', 'currentSemester', 'birthDate', 'careerFinancing', 'school',
                'admissionMode', 'admissionYear'
            ],
            "Docente": [
                'dni', 'firstName', 'lastName', 'address', 'gender', 'phone', 'age', 'birthPlace',
                'currentLiving', 'birthDate', 'workingAtUNAS', 'degree', 'profession', 'maritalStatus'
            ],
            "Administrativo": [
                'dni', 'firstName', 'lastName', 'address', 'gender', 'phone', 'age', 'birthPlace',
                'currentLiving', 'birthDate', 'workingAtUNAS', 'degree', 'profession', 'maritalStatus'
            ],
            "Externo": [
                'dni', 'firstName', 'lastName', 'address', 'gender', 'phone', 'age', 'birthPlace',
                'currentLiving', 'birthDate', 'degree', 'profession', 'maritalStatus'
            ]
        };

        const requiredFields = fieldsByType[this.typePatient] || [];
        this.hasError = requiredFields.some(field => !this.formData[field]);
    }

    previousStep() {
        if (this.currentStep > 1) {
            this.currentStep--;
            this.stepChanged.emit(this.currentStep);
            localStorage.removeItem('formData');
        }
    }
}
