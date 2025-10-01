import {Routes} from "@angular/router";
import {CalendarComponent} from "./pages/calendar/calendar.component";
import {ServicesComponent} from "./pages/services/services.component";
import {DataStatisticComponent} from "./pages/data-statistic/data-statistic.component";
import {StatuteRegulationsComponent} from "./pages/statute-regulations/statute-regulations.component";
import {PostulationComponent} from "./pages/postulation/postulation.component";
import {AnnouncementComponent} from "./pages/announcement/announcement.component";
import {authGuard} from "../../core/guards/auth.guard";
import {UsersComponent} from "./pages/users/users.component";
import {RequestsComponent} from "./pages/requests/requests.component";
import { ThesisPracticingComponent } from "./pages/thesis-practicing/thesis-practicing.component";
import {CalendarSettingComponent} from "./pages/calendar-setting/calendar-setting.component";
import {HouseholdTargetingComponent} from "./pages/household-targeting/household-targeting.component";
import {ScholarshipStudentsComponent} from "./pages/scholarship-students/scholarship-students.component";
import {ApplicationNoticeComponent} from "./pages/application-notice/application-notice.component";
import { VisitaGeneralComponent } from "./pages/visitas/visita-general/visita-general.component";
import { VisitaResidenteComponent } from "./pages/visitas/visita-residente/visita-residente.component";
import { ExamenToxicologicoComponent } from './pages/exam_toxicologico/examen_toxicologico/examen_toxicologico.component';
import { ReporteVisitasComponent } from "./pages/visitas/reporte-visitas/reporte-visitas.component";

import { ResidencesComponent } from "./pages/residences/residences.component";
import { ReportResidencesComponent } from "./pages/report-residences/report-residences.component";
import {PatientsComponent} from "../health-areas/patients/patients.component";
import {NursingComponent} from "../health-areas/nursing/nursing.component";
import {DentistryComponent} from "../health-areas/dentistry/dentistry.component";
import {MedicineComponent} from "../health-areas/medicine/medicine.component";
import { StudentsResidencesComponent } from "./pages/students-residences/students-residences.component";
import {ReportsNursingComponent} from "../health-areas/reports/reports-nursing/reports-nursing.component";
import {ReportsMedicineComponent} from "../health-areas/reports/reports-medicine/reports-medicine.component";
import {ReportsDentistryComponent} from "../health-areas/reports/reports-dentistry/reports-dentistry.component";
import { ResultadoCuestionarioComponent } from "./pages/psicopedagogia/resultado-cuestionario/resultado-cuestionario.component";
import { AtencionesComponent } from "./pages/psicopedagogia/atenciones/atenciones.component";
import { ExpedientesComponent } from "./pages/psicopedagogia/expedientes/expedientes.component";
import { QuestionnairesComponent } from "./pages/psicopedagogia/questionnaires/questionnaires.component";
import { SurveysComponent } from "./pages/psicopedagogia/surveys/surveys.component";
import { CitasPsicologicasComponent } from "./pages/psicopedagogia/citas-psicologicas/citas-psicologicas.component";
import { DataStatisticSrqComponent } from "./pages/psicopedagogia/data-statistic-srq/data-statistic-srq.component";
import {moduleGuard} from "../../core/guards/module-guard.guard";
import {ReportsAdminComponent} from "../health-areas/reports/reports-admin/reports-admin.component";

export const homesRoutes: Routes = [
  {
    path: 'calendar',
    component: CalendarComponent
  },
  {
    path: 'services',
    component: ServicesComponent
  },
  {
    path: 'statistics-data',
    component: DataStatisticComponent,
    canActivate: [authGuard],
    data: {
      role: [1, 2, 3]
    }
  },
  {
    path: 'statues-regulations',
    component: StatuteRegulationsComponent
  },
  {
    path: 'postulation',
    component: PostulationComponent
  },
  {
    path: 'requests',
    component: RequestsComponent,
    canActivate: [authGuard],
    data: {
      role: [1, 2, 3]
    }
  },
  {
    path: 'announcement',
    component: AnnouncementComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'users',
    component: UsersComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'thesis-practicing',
    component: ThesisPracticingComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'calendar-setting',
    component: CalendarSettingComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'residences',
    component: ResidencesComponent,
    canActivate: [authGuard],
    data: {
      role: [1,9]
    }
  },
  {
    path: 'application-notice',
    component: ApplicationNoticeComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }

  },
  {
    path: 'residences/reports',
    component: ReportResidencesComponent,
    canActivate: [authGuard],
    data: {
      role: [1,9]
    }

  },
  {
    path: 'household-targeting',
    component: HouseholdTargetingComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }

  },
  {
    path: 'residences/students',
    component: StudentsResidencesComponent,
    canActivate: [authGuard],
    data: {
      role: [1,9]
    }

  },
  {
    path: 'reports-nursing',
    component: ReportsNursingComponent,
    canActivate: [authGuard],
    data: {
      role: [5]
    }
  },
  {
    path: 'reports-medicine',
    component: ReportsMedicineComponent,
    canActivate: [authGuard],
    data: {
      role: [6]
    }
  },
  {
    path: 'reports-dentistry',
    component: ReportsDentistryComponent,
    canActivate: [authGuard],
    data: {
      role: [7]
    }
  },
  {path: 'patients',
    component: PatientsComponent,
    canActivate: [authGuard],
    data: {
      role: [1, 5, 7]
    },
    children: [
      {
        path: 'create',
        loadComponent: () => import('../health-areas/patients/pages/patient-create/patient-create.component')
      },
      {
        path: ':id/edit',
        loadComponent: () => import('../health-areas/patients/pages/patient-edit/patient-edit.component')
      }
    ]
  },
  {
    path: 'nursing-consultation',
    component: NursingComponent,
    canActivate: [moduleGuard],
    data: {
      module: 'enfermería'
    },
    children: [
      {
        path: 'create',
        loadComponent: () => import('../health-areas/nursing/pages/consultation-create/consultation-create.component')
      },
      {
        path: 'view/paciente/:dni/consulta/:consultationID',
        loadComponent: () => import('../health-areas/nursing/pages/consultation-view/consultation-view.component')
      }
    ]
  },
  {
    path: 'medicine-consultation',
    component: MedicineComponent,
    canActivate: [authGuard],
    data: {
      role: [6]
    },
    children: [
      {
        path: 'edit/paciente/:dni/consulta/:consultationID',
        loadComponent: () => import('../health-areas/medicine/pages/consultation-edit/consultation-edit.component')
      }
    ]
  },
  {
    path: 'dental-consultation',
    component: DentistryComponent,
    canActivate: [authGuard],
    data: {
      role: [7]
    },
    children: [
      {
        path: 'create',
        loadComponent: () => import('../health-areas/dentistry/pages/consultation-create/consultation-create.component')
      },
      {
        path: 'view/paciente/:dni/consulta/:consultationID',
        loadComponent: () => import('../health-areas/dentistry/pages/consultation-edit/consultation-edit.component')
      },
      {
        path: 'edit/nursing/paciente/:dni/consulta/:consultationID',
        loadComponent: () => import('../health-areas/dentistry/pages/consultation-edit-nursing/consultation-edit.component')
      }
    ]
  },
  {
    path: 'admin-reports-medical-areas',
    component: ReportsAdminComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: '',
    redirectTo: 'services',
    pathMatch: 'full'
  },
  //LuisReynaga: routes
  //psicologia
  {
    path: 'resultado',
    canActivate: [authGuard],
    component: ResultadoCuestionarioComponent,
    data: {
      role: [1, 4]
    }
  },
  {
    path: 'atenciones',
    canActivate: [authGuard],
    component: AtencionesComponent,
    data: {
      role: [1, 4]
    }
  },
  {
    path: 'expedientes',
    canActivate: [authGuard],
    component: ExpedientesComponent,
    data: {
      role: [1, 4]
    }
  },
  {
    path: 'cuestionario',
    component: QuestionnairesComponent,
  },
  {
    path: 'encuestas',
    canActivate: [authGuard],
    component: SurveysComponent,
    data: {
      role: [1, 4]
    }
  },
  {
    path: 'citas-psicopedagogicas',
    component: CitasPsicologicasComponent
  },
  {
    path: 'statistics-data-srq',
    component: DataStatisticSrqComponent,
    canActivate: [authGuard],
    data: {
      role: [1, 2, 3]
    }
  },
  //end


  // Ario Masgo: routes
  //visitas
  {
    path: 'visita-general',
    component: VisitaGeneralComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'visita-residente',
    component: VisitaResidenteComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'examen-toxicologico',
    component: ExamenToxicologicoComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  },
  {
    path: 'reporte-visitas',
    component: ReporteVisitasComponent,
    canActivate: [authGuard],
    data: {
      role: [1]
    }
  }

]
