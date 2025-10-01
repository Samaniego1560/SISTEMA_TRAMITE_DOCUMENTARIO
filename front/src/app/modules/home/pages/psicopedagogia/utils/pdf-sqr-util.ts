// utils/pdfGenerator.util.ts
import pdfMake from 'pdfmake/build/pdfmake';
import pdfFonts from 'pdfmake/build/vfs_fonts';

pdfMake.vfs = pdfFonts.vfs;

let isExterno: boolean = false;

export function generateSRQPDF(responseData: any, toastService?: any, base64LogoUnas?: any, base64LogoPsico?: any) {
    if (!responseData) {
        console.error('No hay datos para generar el PDF');
        if (toastService) toastService.add({ type: 'error', message: 'No hay datos para generar el PDF' });
        return;
    }

    const participant = responseData.participant;
    const questions = responseData.questions;
    isExterno = participant.tipo_participante !== 'Alumno';

    const docDefinition: any = {
        pageSize: 'A4',
        pageMargins: [15, 15, 15, 15],
        content: [
            {
                table: {
                    widths: [60, '*', 60],
                    body: [
                        [
                            { image: base64LogoUnas, width: 50, alignment: 'center', margin: [0, 0, 0, 0] },
                            {
                                table: {
                                    widths: ['*', '*', '*', '*'],
                                    body: [
                                        [
                                            {
                                                text: 'UNIVERSIDAD NACIONAL AGRARIA DE LA SELVA',
                                                style: 'headerCell',
                                                fillColor: '#117864',
                                                alignment: 'center',
                                                margin: [2, 6, 2, 6]
                                            },
                                            {
                                                text: 'VICERRECTORADO ACADÉMICO',
                                                style: 'headerCell',
                                                fillColor: '#148F77',
                                                alignment: 'center',
                                                margin: [2, 6, 2, 6]
                                            },
                                            {
                                                text: 'DIRECCIÓN DE BIENESTAR UNIVERSITARIO',
                                                style: 'headerCell',
                                                fillColor: '#17A589',
                                                alignment: 'center',
                                                margin: [2, 6, 2, 6]
                                            },
                                            {
                                                text: 'ÁREA DE PSICOPEDAGOGÍA',
                                                style: 'headerCell',
                                                fillColor: '#1ABB9C',
                                                alignment: 'center',
                                                margin: [2, 6, 2, 6]
                                            }
                                        ]
                                    ],
                                },
                                layout: 'noBorders',
                                margin: [0, 10, 0, 0]

                            },
                            { image: base64LogoPsico, width: 50, alignment: 'center', margin: [0, 0, 0, 0] }
                        ]
                    ]
                },
                layout: 'noBorders',
                margin: [0, 0, 0, 10]
            },
            { text: 'FORMATO Nº 01', style: 'formate', margin: [0, 0, 0, 5] },
            {
                columns: [
                    {
                        text: `Tingo María: ${new Date(participant.created_at).toLocaleDateString('es-PE')}`,
                        alignment: 'left',
                        style: 'info',
                        margin: [20, 0, 0, 0]
                    },
                    {
                        text: 'N°: ____________',
                        alignment: 'right',
                        style: 'info',
                        margin: [0, 0, 20, 0]
                    }
                ],
            },
            { text: 'CUESTIONARIO DE S.R.Q', style: 'title', margin: [0, 8, 0, 2] },
            {
                table: {
                    widths: ['*', '*'],
                    body: [
                        [
                            {
                                stack: [
                                    fieldRow('Apellidos y nombres:', `${participant.nombre} ${participant.apellido}`),
                                    fieldRow('Fecha de nacimiento:', participant.fecha_nacimiento),
                                    fieldRow('Sexo:', participant.sexo),
                                    fieldRow('Lugar de nacimiento:', participant.lugar_nacimiento),
                                    fieldRow('Carrera profesional:', participant.escuela),
                                    fieldRow('Modalidad de ingreso a la UNAS:', participant.modalidad_ingreso),
                                    fieldRow('Año de ingreso a la universidad:', participant.anio_ingreso),
                                    fieldRow('Grado de Instrucción:', participant.grado_instruccion),
                                    fieldRow('Dirección:', participant.direccion),
                                    fieldRow('Labora actualmente en la UNAS:', participant.labora_en_unas),
                                ]
                            },
                            {
                                stack: [
                                    fieldRow('DNI:', participant.dni),
                                    fieldRow('Edad:', participant.edad),
                                    fieldRow('Colegio de procedencia:', participant.colegio_procedencia),
                                    fieldRow('Código de estudiante:', participant.codigo_estudiante),
                                    fieldRow('Con quien vive actualmente:', participant.con_quienes_vive_actualmente),
                                    fieldRow('Semestre que cursa:', participant.semestre_cursa),
                                    fieldRow('N° de celular:', participant.num_telefono),
                                    fieldRow('Estado Civil:', participant.estado_civil),
                                    fieldRow('Profesión:', participant.profesion),

                                ]
                            }
                        ]
                    ]
                },
                layout: {
                    hLineWidth: () => .5,
                    vLineWidth: () => .5,
                },
                margin: [0, 10, 0, 5]
            },
            {
                text: 'RESULTADOS: A continuación se presentan las respuestas seleccionadas por el evaluado.',
                style: 'instructions',
                margin: [0, 0, 0, 10],
            },
            {
                table: {
                    widths: [20, '*', 20, 20, 20, '*', 20, 20],
                    headerRows: 1,
                    body: [
                        [
                            { text: 'N°', style: 'tableHeader' },
                            { text: 'PREGUNTA', style: 'tableHeader' },
                            { text: 'SÍ', style: 'tableHeader' },
                            { text: 'NO', style: 'tableHeader' },
                            { text: 'N°', style: 'tableHeader' },
                            { text: 'PREGUNTA', style: 'tableHeader' },
                            { text: 'SÍ', style: 'tableHeader' },
                            { text: 'NO', style: 'tableHeader' }
                        ],
                        ...generateQuestionRows(questions)
                    ]
                },
                layout: {
                    fillColor: (rowIndex: number) => (rowIndex === 0 ? '#1ABB9C' : null),
                    hLineWidth: () => .5,
                    vLineWidth: () => .5,
                }
            },
            {
                stack: [
                    {
                        columns: [
                            {
                                width: '40%',
                                stack: [
                                    { text: '________________________________', alignment: 'center', margin: [0, 40, 0, 0] },
                                    { text: 'FIRMA', style: 'signature', alignment: 'center', margin: [0, 5, 0, 20] }
                                ]
                            },
                            { width: '*', text: '' }, // espacio flexible
                            {
                                width: '40%',
                                columns: [
                                    {
                                        width: '*',
                                        stack: [
                                            { text: 'Síguenos como:', style: 'info', alignment: 'right', margin: [0, 40, 20, 2] },
                                            {
                                                text: 'https://web.facebook.com/psicopedagogiaunas',
                                                link: 'https://web.facebook.com/psicopedagogiaunas',
                                                color: 'blue',
                                                style: 'info',
                                                alignment: 'right',
                                                margin: [0, 0, 20, 0]
                                            }
                                        ]
                                    },
                                    {
                                        width: 'auto',
                                        qr: 'https://web.facebook.com/psicopedagogiaunas',
                                        fit: 60,
                                        alignment: 'right',
                                        margin: [10, 40, 10, 0]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }

        ],
        styles: {
            headerCell: { fontSize: 8, alignment: 'center', color: 'black', margin: [2, 2, 2, 2] },
            title: { fontSize: 12, alignment: 'center', bold: true },
            formate: { fontSize: 9, alignment: 'center', bold: true },
            info: { fontSize: 8, margin: [2, 1, 2, 1] },
            instructions: { fontSize: 9 },
            tableHeader: { fontSize: 9, alignment: 'center', color: 'black', margin: [2, 2, 2, 2] },
            tableCell: { fontSize: 8, alignment: 'center', margin: [2, 2, 2, 2] },
            questionText: { fontSize: 8, alignment: 'left', margin: [2, 2, 2, 2] },
            signature: { fontSize: 7, alignment: 'center' }
        }
    };

    if (isExterno) {
        docDefinition.content.splice(4, 0, { text: '(EXTERNOS)', alignment: 'center', bold: true, fontSize: 9});
    }

    try {
        const code = participant.tipo_participante === 'Alumno' ? participant.codigo_estudiante : participant.dni;
        const namePdf = `${participant.nombre} ${participant.apellido}-${code}`
        pdfMake.createPdf(docDefinition).open();
        // pdfMake.createPdf(docDefinition).download(`${namePdf}.pdf`);
        console.log('PDF generado exitosamente');
    } catch (error) {
        console.error('Error al generar PDF:', error);
        if (toastService) toastService.add({ type: 'error', message: 'Error al generar el PDF' });
    }
}

export function getBase64ImageFromUrl(url: string): Promise<string> {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.setAttribute('crossOrigin', 'anonymous');
        img.src = url;
        img.onload = () => {
            const canvas = document.createElement('canvas');
            canvas.width = img.width;
            canvas.height = img.height;
            const ctx = canvas.getContext('2d');
            if (ctx) {
                ctx.drawImage(img, 0, 0);
                const dataURL = canvas.toDataURL('image/png');
                resolve(dataURL);
            } else {
                reject('No se pudo obtener el contexto del canvas');
            }
        };
        img.onerror = () => reject('No se pudo cargar la imagen: ' + url);
    });
}



function fieldRow(label: any, value: any) {
    const labelNotExterno = ['Carrera profesional:', 'Modalidad de ingreso a la UNAS:', 'Año de ingreso a la universidad:', 
        'Colegio de procedencia:', 'Código de estudiante:', 'Semestre que cursa:'
    ];
    const labelNotStudent =  ['Grado de Instrucción:', 'Dirección:', 'Estado Civil:', 'Profesión:', 'Labora actualmente en la UNAS:']; 
    console.log(isExterno)
    if(isExterno && labelNotExterno.includes(label)){
        return
    }
    if(!isExterno && labelNotStudent.includes(label)){
        return
    }
    const type = typeof value;
    if(type == 'boolean' && value === true){
        value = 'Sí'
    }else if(type == 'boolean' && value === false){
        value = 'No'
    }
    return {
        columns: [
            { text: label, style: 'info', width: '50%', bold: true },
            { text: value || '', style: 'info', alignment: 'left' }
        ]
    };
}

function generateQuestionRows(questions: any[]): any[] {
    const rows = [];

    // Asegúrate de que hay al menos 30 preguntas
    const leftQuestions = questions.slice(0, 15);
    const rightQuestions = questions.slice(15, 30);

    for (let i = 0; i < 15; i++) {
        const left = leftQuestions[i] || {};
        const right = rightQuestions[i] || {};

        const row = [
            { text: left.Number?.toString() || '', style: 'tableCell', },
            { text: left.Text || '', style: 'questionText' },
            { text: left.Answer === 'Si' ? 'X' : '', style: 'tableCell' },
            { text: left.Answer === 'No' ? 'X' : '', style: 'tableCell' },

            { text: right.Number?.toString() || '', style: 'tableCell' },
            { text: right.Text || '', style: 'questionText' },
            { text: right.Answer === 'Si' ? 'X' : '', style: 'tableCell' },
            { text: right.Answer === 'No' ? 'X' : '', style: 'tableCell' }
        ];

        rows.push(row);
    }

    return rows;
}


