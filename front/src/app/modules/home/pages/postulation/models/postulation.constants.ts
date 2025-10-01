import {SocioEconomicSheet} from "./postulation.models";

export const structureSocioEconomicSheet: SocioEconomicSheet = {
  sections: [
    {
      name: 'DATOS PERSONALES',
      type_section: 'form',
      attributes: [
        {name: 'nombres', description: 'Nombres', type_attribute: 'input', value: '', order: 12},
        {name: 'date_age', description: 'Fecha de Nacimiento', type_attribute: 'input', value: '', order: 3},
        {name: 'edad', description: 'Edad', type_attribute: 'input', value: '', order: 3},
        {name: 'sexo', description: 'Sexo', type_attribute: 'select', value: '', options: 'M,F', order: 3},
        {name: 'dni', description: 'Dni', type_attribute: 'input', value: '', order: 3},
        {name: 'marital_status', description: 'Estado civil', type_attribute: 'input', value: '', order: 12},
        {name: 'address_origin', description: 'Domicilio de Procedencia', type_attribute: 'input', value: '', order: 12},
        {name: 'district', description: 'Distrito', type_attribute: 'input', value: '', order: 4},
        {name: 'province', description: 'Provincia', type_attribute: 'input', value: '', order: 4},
        {name: 'department', description: 'Departamento', type_attribute: 'input', value: '', order: 4},
        {name: 'personal_cell_phone', description: 'Celular personal', type_attribute: 'input', value: '', order: 6},
        {name: 'cell_phone_proxy', description: 'Celular apoderado', type_attribute: 'input', value: '', order: 6},
      ]
    },
    {
      name: 'DATOS ACADÉMICOS',
      type_section: 'form',
      attributes: [
        {name: 'facultad', description: 'Facultad', type_attribute: 'input', value: '', order: 4},
        {name: 'escuelaprofesional', description: 'Escuela', type_attribute: 'input', value: '', order: 4},
        {name: 'códigoestudiante', description: 'Codigo', type_attribute: 'input', value: '', order: 4},
        {name: 'modalidaddeingreso', description: 'Modalidad de Ingreso a la UNAS', type_attribute: 'radio', value: '',
          options: 'Admisión Ordinaria,Primeros Puestos,Segunda Profesión,CEPREUNAS,Traslado Externo,Traslado Interno,Por Discapacidad, Victimas por Terrorismo,Convenio, Deportista Calificado,Arte y Cultura,Otros'
          , order: 12},
        {name: 'which_school', description: 'En qué Colegio Estudio', type_attribute: 'input', value: '', order: 12},
        {name: 'school_type', description: 'Tipo de Colegio', type_attribute: 'radio', value: '', options: 'Estatal,Particular,Parroquial,Paraestatal', order: 12},
      ]
    },
    {
      name: 'COMPOSICIÓN FAMILIAR (NUCLEO DE CONVIVENCIA)',
      type_section: 'table',
      attributes: [
        {name: 'family_name_lastname', description: 'Nombres y Apellidos', type_attribute: 'input', value: '', order: 6},
        {name: 'family_age', description: 'Edad', type_attribute: 'input', value: '', order: 6},
        {name: 'family_marital_status', description: 'Estado civil', type_attribute: 'input', value: '', order: 6},
        {name: 'family_educational_level', description: 'Nivel educativo', type_attribute: 'input', value: '', order: 6},
        {name: 'family_occupation', description: 'Ocupación', type_attribute: 'input', value: '', order: 12},
      ]
    },
    {
      name: 'INGRESO/DEPENDENCIA ECONÓMICA',
      type_section: 'form',
      attributes: [
        {name: 'contribute_financially', description: ' ¿Quién o quiénes aportan económicamente en tu hogar?', type_attribute: 'radio', options: 'Papá,Mamá,Hermanos,El Estudiante,Otros', value: '', order: 12},
        {name: 'monthly_economic_income', description: '¿Cuánto es el ingreso económico mensual de tu hogar?', type_attribute: 'input', value: '', order: 12},
      ]
    },
    {
      name: 'VIVIENDA',
      type_section: 'form',
      attributes: [
        {name: 'tenure', description: 'Tenencia', type_attribute: 'select', options: 'Propia,Hipotecado,Alquilada,Alojado,Guardiania,Invasión', value: '', order: 12},
        {name: 'wall', description: 'Pared', type_attribute: 'select', options: 'Noble,Adobe,Madera,Quincha,Esteras,Otros',  value: '', order: 4},
        {name: 'floor', description: 'Piso', type_attribute: 'select', options: 'Tierra,Parquet,Céramica,Concreto,Cemento,Otros',  value: '', order: 4},
        {name: 'ceiling', description: 'Techo', type_attribute: 'select', options: 'Calamina,Esteras,Paja,Madera,Plastico,Otros',  value: '', order: 4},
        {name: 'housing_area', description: 'Zona de la Vivienda', type_attribute: 'select', options: 'Urbana,Rural,Residencial,AAHH,Otros',  value: '', order: 6},
        {name: 'equipment_housing', description: 'Equipamiento de la Vivienda', type_attribute: 'select', options: 'Equipo de sonido,Refrigeradora/Congeladora,Telvisor,Cocina a Gas,Licuadora,Lavadora,Computadora,HornoMicrondas',  value: '', order: 6},
      ]
    },
    {
      name: 'ACCESO A LOS SERVICIOS BÁSICOS',
      type_section: 'form',
      attributes: [
        {name: 'account_services_1', description: ' ¿Con qué servicios cuentas?', type_attribute: 'checkbox', options: 'Agua potable,Desagüe,energía eléctrica', value: '', order: 12},
        {name: 'account_services_2', description: ' ¿Con qué servicios cuentas?', type_attribute: 'checkbox', options: 'Internet,Cable,Otros', value: '', order: 12},
        {name: 'linear_motorcycle', description: ' ¿Tiene Moto lineal?', type_attribute: 'radio', options: 'Si,No', value: '', order: 12},
        {name: 'others_services', description: ' ¿Que otra tenencial Tiene?', type_attribute: 'input', value: '', order: 12},
      ]
    },
    {
      name: 'SALUD DEL ALUMNO',
      type_section: 'form',
      attributes: [
        {name: 'suffer_illness', description: '¿Padece alguna enfermedad?', type_attribute: 'radio', options: 'Si,No', value: '', order: 12},
        {name: 'has_treatment', description: '¿Tiene tratamiento?', type_attribute: 'radio', options: 'Si,No', value: '', order: 12},
        {name: 'where_serve', description: '¿Donde se Atiende?', type_attribute: 'radio', options: 'MINSA SIS,Servicio Médico OBU,ESSALUD,Clínica,FOSPOLI,Ninguno', value: '', order: 12},
        {name: 'allergic_medication', description: 'Es alérgica(o) a algún medicamento', type_attribute: 'radio', options: 'Si,No', value: '', order: 12},
        {name: 'blood_group', description: 'Cuál es su grupo sanguíneo', type_attribute: 'input', options: '', value: '', order: 12},
      ]
    },
    {
      name: 'FOTOGRAFIA DE FACHADA DE SU HOGAR Y REFERENCIA (ADJUNTAR)',
      type_section: 'form',
      attributes: [
        {name: 'house_facade', description: 'Fachada Casa', type_attribute: 'file-image', options: 'png,jpg,jpeg', value: '', order: 12},
      ]
    },

  ]
}
