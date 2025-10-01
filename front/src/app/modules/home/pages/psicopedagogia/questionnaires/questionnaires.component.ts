import { Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NgFor, NgIf } from "@angular/common";
import { Subscription } from "rxjs";
import { PreQuestionsComponent } from "../pre-questions/pre-questions.component";
import { ActivatedRoute } from '@angular/router';
import { HeaderComponent } from '../../../../../core/ui/header/header.component';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { BlockUiComponent } from '../../../../../core/ui/block-ui/block-ui.component';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { UserAgent } from '../../../../../core/utils/functions/userAnget';
import { DataInit } from '../interfaces/data-init';


@Component({
    selector: 'app-questionnaires',
    standalone: true,
    imports: [
        ReactiveFormsModule,
        HeaderComponent,
        PreQuestionsComponent,
        ToastComponent,
        NgIf,
        NgFor,
        BlockUiComponent,
        FormsModule
    ],
    templateUrl: './questionnaires.component.html',
    styleUrl: './questionnaires.component.scss',
    providers: [ToastService, ManagerService]
})
export class QuestionnairesComponent implements OnInit {

    private _subscriptions: Subscription = new Subscription();
    protected openMenu: boolean = !UserAgent.IsMobileDevice();
    showFormSRQ: boolean = false;
    protected isLoad: boolean = false;
    questions: any = [];
    responses: { [key: number]: string } = {};
    fullName: string = '';
    currentSemester!: string;
    schoolOrigin!: string;
    universityYear!: string;
    phone!: string;
    livingWith!: string;
    questionId: Number = -1;
    dataStudent: any = null;
    isSend: boolean = false;
    negativeResult: boolean = false;
    direccion!: string
    idEncuesta!: number;
    formData: any;
    isGenerateSurvery = false;
    hideForm = false;
    key_url = ''
    byPassValidation = false;
    currentDate = new Date().getFullYear();
    typeUser: string = '';

    degree!: string;
    profession!: string;
    gender!: string;
    age: number = 0;
    birthPlace!: string;
    birthDate!: string;
    maritalStatus!: string;

    ciclos: string[] = [
        "1°", "2°", "3°", "4°", "5°",
        "6°", "7°", "8°", "9°", "10°"
    ];

    constructor(
        private _managerService: ManagerService,
        private _toastService: ToastService,
        private route: ActivatedRoute
    ) { }

    ngOnInit() {
        this.route.queryParams.subscribe(params => {
            if (params['dataId']) {
                try {
                    const decodedData = decodeURIComponent(params['dataId']);
                    this.hideForm = true;
                    const jsonString = atob(decodedData);
                    this.formData = JSON.parse(JSON.parse(jsonString));
                    this.key_url = this.formData.key_url;
                    this.byPassValidation = true
                    this.keyUrlExists();
                } catch (error: any) {
                    console.error("Error al decodificar los datos:", error);
                }
            }
        });
    }

    onSubmit(): void {
        if (this.questions.length != Object.values(this.responses).length) {
            this._toastService.add({ type: 'warning', message: 'Por favor, completa todas las preguntas S.R.Q' });
            return;
        }
        const answers = Object.entries(this.responses).map(([key, value]) => {
            return { id_pregunta: Number(key), respuesta: value };
        });

        const isNotValidStudent = (this.phone === undefined || this.schoolOrigin === undefined || this.universityYear === undefined ||this.currentSemester === undefined ||
                        this.livingWith === undefined);
        const isNotValidTesista = (this.degree === undefined || this.profession === undefined || this.gender === undefined || this.age === 0 || 
                                    this.birthPlace === undefined || this.birthDate === undefined || this.maritalStatus === undefined)
        if(!this.byPassValidation){
            if (isNotValidStudent && this.typeUser === 'Alumno') {
                this._toastService.add({ type: 'warning', message: 'Por favor, completa todos los campos antes de enviar.' });
                return;
            }
            if (isNotValidTesista && this.typeUser === 'Tesista') {
                this._toastService.add({ type: 'warning', message: 'Por favor, completa todos los campos antes de enviar.' });
                return;
            }
        }


        if (this.phone === undefined || !this.phone.match(/^9\d{8}$/)) {
            this._toastService.add({
                type: 'warning',
                message: 'Por favor, ingrese un número de teléfono válido.'
            });
            return;
        }
        let alumno = this.dataStudent.detalle;
        if(this.typeUser === 'Tesista'){
            const detail = {
                num_telefono: this.phone,
                grado_instruccion: this.degree,
                profesion: this.profession,
                sexo: this.gender,
                edad: Number(this.age),
                lugar_nacimiento: this.birthPlace,
                fecha_nacimiento: this.birthDate,
                estado_civil: this.maritalStatus,
            }
            alumno = {...alumno, ...detail};
        }
        const payload = {
            participante: {
                tipo_participante: alumno.tipo_participante || "Alumno",
                nombre: alumno.nombre || alumno.nombres,
                apellido: alumno.apellido || `${alumno.apellido_paterno} ${alumno.apellido_materno}`,
                dni: alumno.dni,
                num_telefono: this.phone,
                con_quienes_vive_actualmente: this.livingWith || "",
                fecha_nacimiento: alumno.fecha_nacimiento,
                edad: alumno.edad,
                lugar_nacimiento: alumno.lugar_nacimiento || this.birthPlace,
                colegio_procedencia: this.schoolOrigin || null,
                anio_ingreso: Number(this.universityYear) || null,
                semestre_cursa: this.currentSemester || null,
                escuela: alumno.escuela || alumno.escuela_profesional,
                codigo_estudiante: alumno.codigo_estudiante,
                modalidad_ingreso: alumno.modalidad_ingreso,
                sexo: alumno.sexo,
                es_srq: true,
                direccion: this.direccion || alumno.direccion,
                quien_financia_carrera: alumno.quien_financia_carrera,
                motivo_consulta: alumno.motivo_consulta,
                diagnostico_id: Number(alumno.diagnostico_id),
                situacion_actual: alumno.situacion_actual,
                otros_procedimientos: alumno.otros_procedimientos,
                grado_instruccion: alumno.grado_instruccion,
                profesion: alumno.profesion,
                estado_civil: alumno.estado_civil,
                key_url: this.key_url
            },
            encuesta_id: this.idEncuesta,
            respuestas: answers || []
        };
        this.isLoad = true;
        this._managerService.saveAnswers(payload).subscribe({
            next: (response: any) => {
                this.isLoad = false;
                this.isSend = true;
                this._toastService.add({ type: 'success', message: response.msg });
                this.negativeResult = response.detalle != 'Aprobado'
            },
            error: (error: any) => {
                console.error('Error al obtener preguntas:', error);
                this.isLoad = false;
                this._toastService.add({ type: 'error', message: 'A ocurrido un error guardando las respuestas' });
            }
        })
        console.log('Form Data:', JSON.stringify(payload));

    }

    public processParameterData() {
        this.showFormSRQ = true;
        this.dataStudent.detalle = this.dataStudent.detalle || {}
        if (this.fullName === '') {
            this.fullName = `${this.formData.lastName}, ${this.formData.firstName}`;
        }
        this.phone = this.formData.phone;
        this.livingWith = this.formData.currentLiving;
        this.schoolOrigin = this.formData.school;
        this.universityYear = this.formData.admissionYear;
        this.currentSemester = this.formData.currentSemester;
        this.dataStudent.detalle.nombre = this.formData.firstName;
        this.dataStudent.detalle.apellido = this.formData.lastName;
        this.dataStudent.detalle.dni = this.formData.dni;
        this.dataStudent.detalle.tipo_participante = this.formData.tipoParticipante;
        this.dataStudent.detalle.edad = Number(this.formData.age);
        this.dataStudent.detalle.lugar_nacimiento = this.formData.birthPlace;
    }

    addIdEncuesta(id: number) {
        console.log('ID Encuesta recibido:', id);
        this.idEncuesta = id;
    }

    public getDataQuestions(dataInit: DataInit): void {
        const { numDocumento, typeUser, fistName, lastName} = dataInit; 
        if (!numDocumento) {
            this._toastService.add({ type: 'error', message: 'El DNI es requerido' });
            return;
        }
        this.isLoad = true;
        this._subscriptions.add(
            this._managerService.getStudentQuestions(numDocumento, typeUser).subscribe({
                next: (response: any) => {
                    this.isLoad = false;
                    if (!response) {
                        this._toastService.add({ type: 'warning', message: 'No se encontraron preguntas para este estudiante' });
                        return;
                    }
                    this.dataStudent = !this.dataStudent ? response[0]: this.dataStudent;
                    this.questions = response[1].detalle;
                    this.questions.sort((a: any, b: any) => {
                        return parseInt(a.order) - parseInt(b.order);
                    });
                    this.questionId = this.questions[0].id;
                    if (this.dataStudent.detalle && typeUser === 'Alumno') {
                        this.showFormSRQ = true;
                        if (this.dataStudent.detalle.apellido_paterno && this.dataStudent.detalle.apellido_materno && this.dataStudent.detalle.nombres) {
                            this.fullName = `${this.dataStudent.detalle.apellido_paterno} ${this.dataStudent.detalle.apellido_materno} ${this.dataStudent.detalle.nombres}`;
                        } else {
                            this.fullName = `${this.dataStudent.detalle.apellido} ${this.dataStudent.detalle.nombre}`;
                        }
                    }else if(typeUser !== 'Alumno'){
                        this.fullName = `${fistName?.toUpperCase()} ${lastName?.toUpperCase()}`;
                        this.showFormSRQ = true;
                    }
                    if (this.hideForm) {
                        this.processParameterData();
                    }

                },
                error: (error: any) => {
                    console.error('Error al obtener preguntas:', error);
                    this.isLoad = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener las preguntas del estudiante' });
                }
            })
        );
    }

    public fetchStudentData(dataInit: DataInit): void {
        const { numDocumento, typeUser, lastName, fistName} = dataInit;
        this.typeUser = typeUser;
        if (!numDocumento) {
            this._toastService.add({ type: 'error', message: 'El DNI es requerido' });
            return;
        }
    
        this.isLoad = true;
        if(typeUser === 'Alumno'){
            this._subscriptions.add(
                this._managerService.getDatosAlumnoAcademico(numDocumento).subscribe({
                    next: (response: any) => {
                        this.isLoad = false;
                        if (response?.detalle?.length == 0 || !response) {
                            this._toastService.add({ type: 'warning', message: 'No se pudo recuperar los datos del estudiante' });
                            return;
                        }
                        this.dataStudent = response;
                        const detail = this.dataStudent.detalle;
                        detail.apellido = `${this.dataStudent.detalle.appaterno} ${this.dataStudent.detalle.apmaterno}`
                        detail.codigo_estudiante = detail.codalumno?.toString(); 
                        detail.dni = numDocumento; 
                        detail.escuela = detail.nomesp; 
                        detail.lugar_nacimiento = this.birthPlace,
                        detail.fecha_nacimiento = detail.fecnac; 
                        detail.modalidad_ingreso = detail.mod_ingreso; 
                        detail.edad = this.calcularEdad(detail.fecnac);
                        this.dataStudent.detalle = detail;
                        if (this.dataStudent) {
                            this.getDataQuestions(dataInit);
                        }
                    },
                    error: (error: any) => {
                        console.error('Error al obtener los datos del estudiante:', error);
                        this.isLoad = false;
                        this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener los datos del estudiante' });
                    }
                })
            );
        }else{
            this.dataStudent = {};
            const detail = {
                dni: numDocumento,
                nombre: fistName,
                apellido: lastName,
                tipo_participante: typeUser,
            }
            this.dataStudent.detalle = detail;
            this.getDataQuestions(dataInit);
        }
        
    }

    public calcularEdad(fechaNacimiento: string): number {
        const fechaNac = new Date(fechaNacimiento);
        const hoy = new Date();
        let edad = hoy.getFullYear() - fechaNac.getFullYear();
        const mes = hoy.getMonth() - fechaNac.getMonth();
        const dia = hoy.getDate() - fechaNac.getDate();
        if (mes < 0 || (mes === 0 && dia < 0)) {
            edad--;
        }
    
        return edad;
    }
    

    public toggleMenu(): void {
        this.openMenu = !this.openMenu;
    }

    public keyUrlExists() {
        this.isLoad = true;
        this._subscriptions.add(
            this._managerService.keyUrlExists(this.formData.key_url).subscribe({
                next: (response: any) => {
                    if(response.detalle == true){
                        this.isLoad = false;
                        this._toastService.add({ type: 'warning', message: 'Usted ya completó la encuesta, gracias por su participación' });
                    }else{
                        const dataInit: DataInit = {
                            numDocumento: this.formData.dni,
                            typeUser: this.formData.tipoParticipante
                        }
                        this.getDataQuestions(dataInit);
                    }
                },
                error: (error: any) => {
                    console.error('Error al verificar data:', error);
                    this.isLoad = false;
                    this._toastService.add({ type: 'error', message: 'Ocurrió un error al obtener las preguntas del estudiante' });
                }
            })
        );
    }
}
