import {ChartOptionsBarra} from "../data-statistic.component";

export const ChartBarsFaculty: ChartOptionsBarra = {
  series: [],
  chart: {
    type: "bar",
    height: 350,
    width: 500
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: "55%"
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    show: true,
    width: 2,
    colors: ["transparent"]
  },
  xaxis: {
    categories: [
      'AGRONOMIA',
      'CIENCIAS CONTABLES',
      'F.C.A',
      'F.I.I.A',
      'F.I.I.S',
      'F.I.M.E',
      'F.R.N.R',
      'ZOOTECNIA'
    ]
  },
  yaxis: {},
  fill: {
    opacity: 1
  },
  tooltip: {},
  legend: {
    position: "bottom",
    horizontalAlign: "left",
    offsetX: 40
  },
  title: {
    align: "center",
    text: 'Beneficiarios por Facultad Profesional'
  }
};

export const ChartBarsSchool: ChartOptionsBarra = {
  series: [],
  chart: {
    type: "bar",
    height: 350,
    width: 500
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: "55%"
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    show: true,
    width: 2,
    colors: ["transparent"]
  },
  xaxis: {
    categories: [
      "ADMINISTRACION",
      "AGRONOMIA",
      "CONTABILIDAD",
      "ECONOMIA",
      "E.I.A",
      "E.I.C.S.A",
      "E.I.I.A",
      "E.I.I.S",
      "E.I.R.N.R",
      "E.I.F",
      "E.I.M.E",
      "ZOOTECNIA"
    ]
  },
  yaxis: {},
  fill: {
    opacity: 1
  },
  tooltip: {},
  legend: {
    position: "bottom",
    horizontalAlign: "left",
    offsetX: 40
  },
  title: {
    align: "center",
    text: 'Beneficiarios por Escuela Profesional'
  }
};

export const ChartBarsSchoolGender: ChartOptionsBarra = {
  series: [],
  chart: {
    type: "bar",
    height: 350,
    width: 500
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: "55%"
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    show: true,
    width: 2,
    colors: ["transparent"]
  },
  xaxis: {
    categories: [
      "ADMINISTRACION",
      "AGRONOMIA",
      "CONTABILIDAD",
      "ECONOMIA",
      "E.I.A",
      "E.I.C.S.A",
      "E.I.I.A",
      "E.I.I.S",
      "E.I.R.N.R",
      "E.I.F",
      "E.I.M.E",
      "ZOOTECNIA"
    ]
  },
  yaxis: {},
  fill: {
    opacity: 1
  },
  tooltip: {},
  legend: {
    position: "bottom",
    horizontalAlign: "left",
    offsetX: 40
  },
  title: {
    align: "center",
    text: 'Beneficiarios de Escuela Profesional por genero'
  }
};
