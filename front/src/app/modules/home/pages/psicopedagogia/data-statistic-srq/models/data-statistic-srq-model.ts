import {
  ApexAxisChartSeries,
  ApexChart,
  ApexDataLabels, ApexFill,
  ApexLegend,
  ApexNonAxisChartSeries, ApexResponsive, ApexStroke,
  ApexTitleSubtitle,
  ApexTooltip, ApexXAxis, ApexYAxis
} from "ng-apexcharts";


export type ChartOptionsPieSqr = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  labels: string[];
  title: ApexTitleSubtitle;
  legend: ApexLegend;
  dataLabels: ApexDataLabels;
  tooltip: ApexTooltip;
  fill: ApexFill;
  responsive: ApexResponsive[];
};

export type ChartOptionsAreaSqr = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  stroke: ApexStroke;
  dataLabels: ApexDataLabels;
  yaxis: ApexYAxis;
  title: ApexTitleSubtitle;
  labels: string[];
  legend: ApexLegend;
  subtitle: ApexTitleSubtitle;
};

export const chartPieSqr: ChartOptionsPieSqr = {
  series: [],
  chart: {
    type: 'pie',
    width: 380
  },
  labels: [],
  title: {
    text: 'Estados de evaluación',
    align: 'center'
  },
  legend: {
    position: 'bottom'
  },
  dataLabels: {
    enabled: true,
    formatter: (val, opts) => {
      const label = opts.w.globals.labels[opts.seriesIndex];
      return `${label}: ${(val as number).toFixed(1)}%`;
    }
  },
  tooltip: {
    enabled: true,
    y: {
      formatter: (val) => `${val} unidades`
    }
  },
  fill: {
    colors: ['#1E90FF', '#00C49A', '#FFBB28', '#FF8042', '#8884D8']
  },
  responsive: [
    {
      breakpoint: 480,
      options: {
        chart: {
          width: 300
        },
        legend: {
          position: 'bottom'
        }
      }
    }
  ]
};
export const chatAreaSqr : ChartOptionsAreaSqr = {
  series: [
    {
      name: "STOCK ABC",
      data: []
    },
    {
      name: "STOCK ABC",
      data: []
    }
  ],
  chart: {
    type: "area",
    height: 350,
    zoom: {
      enabled: false
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    curve: "straight"
  },

  title: {
    text: "Participantes del mes",
    align: "left"
  },
  subtitle: {
    text: "",
    align: "left"
  },
  labels: [],
  xaxis: {
    type: "datetime"
  },
  yaxis: {
    opposite: true
  },
  legend: {
    horizontalAlign: "left"
  }
};

export const trimestres = [
  {name: '1er trimestre', value: '1,2,3'},
  {name: '2do trimestre', value: '4,5,6'},
  {name: '3er trimestre', value: '7,8,9'},
  {name: '4to trimestre', value: '10,11,12'},
]

export const cuadros = [
  {name: '1. CONSULTAS Y ATENCIONES PSICOLOGICAS A ESTUDIANTES', value: '1'},
  {name: '2. CONSULTAS Y ATENCIONES PSICOLOGICAS A DOCENTES, NO DOCENTES Y EXTERNOS', value: '2'},
  {name: '3. CONSULTAS DE RESULTADOS DEL CUESTIONARIO S.R.Q', value: '3'},

]
