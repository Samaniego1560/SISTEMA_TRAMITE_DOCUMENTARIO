import {
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  inject,
  Input,
  OnChanges,
  output,
  signal,
  SimpleChanges,
  viewChild
} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule,} from "@angular/forms";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {FillColorDirective} from "../../../../../core/directives/fill-color.directive";
import {OdontogramExam} from "../../../../../core/models/areas/dentistry.model";
import {FillColorDefaultDirective} from "../../../../../core/directives/fill-color-default.directive";
import {odontogramConstant} from "../../constants/odontogram.constant";

@Component({
  selector: 'app-oral-exam-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective,
    FormsModule,
    FillColorDirective,
    FillColorDefaultDirective
  ],
  templateUrl: './oral-exam-form.component.html',
  styleUrl: './oral-exam-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class OralExamFormComponent extends FormBaseComponent<OdontogramExam> implements OnChanges {

  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  screenOdontogram = viewChild<ElementRef<HTMLDivElement>>('screenOdontogram');
  svgOdontogram = viewChild<ElementRef<SVGElement>>('svgOdontogram');

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('examen_bucal');
  public shadowColor = signal('none');
  colors: Record<string, string> = {
    none: '#ffffff',
    blue: '#004da1',
    red: '#f80102',
    green: '#008f39',
    yellow: '#ffff00',
    wine: '#5e2129 '
  }
  colorsPath: Record<string, string> = {};

  private _fb = inject(FormBuilder);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoData']?.currentValue && this.infoData) {
      this.form.patchValue(this.infoData);
      this.refreshColorsOdontogram(this.getSvgFromBase64(this.infoData.odontograma_img));
      this.active.set(false);
      this.form.disable()
    }
  }

  refreshColorsOdontogram(svgString: string) {
    if(svgString === '') return

    const parser = new DOMParser()
    const svg = parser.parseFromString(svgString, 'image/svg+xml')
    const paths = svg.querySelectorAll('path')
    for (const path of Array.from(paths)) {
      if (!path.id || !path.style.fill) continue;
      this.colorsPath[path.id] = path.style.fill
    }
  }

  private _buildForm() {
    return this._fb.group({
      cpod: ['',],
      ihos: ['',],
      observacion: ['',],
      comentarios: ['',],
      odontograma_img: ['',],
    })
  }

  public get name() {
    return this._name();
  }

  public isValid(): boolean {
    return this.form.valid
  }

  public getData() {
    if(this.active()) {
      this.refreshColorsOdontogram(this.getSvgAsString());
    }

    const svg = this.getBase64Odontogram();
    return {
      ...this.form.value,
      odontograma_img: svg
    }
  }

  getSvgAsString(): string {
    const svgCopy = this.svgOdontogram()?.nativeElement.cloneNode(true) as SVGElement;

    if (!svgCopy.hasAttribute('xmlns')) {
      svgCopy.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    }

    this.cleanAngularAttributes(svgCopy);

    return new XMLSerializer().serializeToString(svgCopy);
  }

  getSvgFromBase64(base64: string): string {
    return decodeURIComponent(escape(atob(base64)));
  }

  getBase64Odontogram() {
    const parser = new DOMParser()
    const svgDoc = parser.parseFromString(odontogramConstant, 'image/svg+xml')
    const svgElement = svgDoc.querySelector('svg')!;
    for (const colorPath in this.colorsPath) {
      const color = this.colorsPath[colorPath];
      const pathElement = svgElement.querySelector(`path[id="${colorPath}"]`)!;
      pathElement.setAttribute('style', `fill: ${color};`);
    }

    const svgString = new XMLSerializer().serializeToString(svgElement);
    return btoa(unescape(encodeURIComponent(svgString)));
  }

  setActiveExam(isActive: boolean): void {
    if (!isActive) {
      this.refreshColorsOdontogram(this.getSvgAsString());
    }
    this.active.set(isActive);
  }

  public setFullScreen() {
    this.screenOdontogram()?.nativeElement.requestFullscreen().then();
  }

  public setFillColor(paths: HTMLElement[]) {
    paths.forEach(path => {
      path.style.fill = this.colors[this.shadowColor()];
    })
  }

  downloadSvg(fileName: string = 'download'): void {
    const svgCopy = this.svgOdontogram()?.nativeElement.cloneNode(true) as SVGElement;

    if (!svgCopy.hasAttribute('xmlns')) {
      svgCopy.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    }

    this.cleanAngularAttributes(svgCopy);

    const svgString = new XMLSerializer().serializeToString(svgCopy);

    const blob = new Blob([svgString], {type: 'image/svg+xml;charset=utf-8'});
    const url = URL.createObjectURL(blob);

    const downloadLink = document.createElement('a');
    downloadLink.href = url;
    downloadLink.download = `${fileName}.svg`;

    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
    URL.revokeObjectURL(url);
  }

  /**
   * Limpia atributos de Angular recursivamente
   */
  private cleanAngularAttributes(element: Element): void {
    const attributesAllowed = [
      'xlmns',
      'fill',
      'id',
      'viewBox',
      'width',
      'class',
      'stroke-linejoin',
      'stroke-width',
      'stroke',
      'd',
      'clip-rule',
      'fill-rule',
      'style'
    ];

    for (let i = element.attributes.length - 1; i >= 0; i--) {
      const attr = element.attributes[i];
      if (!attributesAllowed.includes(attr.name)) {
        element.removeAttribute(attr.name);
      }
    }

    Array.from(element.children).forEach(child => this.cleanAngularAttributes(child));
  }
}
