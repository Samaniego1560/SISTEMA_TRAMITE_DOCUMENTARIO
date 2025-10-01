import {
  Component,
  computed, effect,
  EventEmitter,
  input,
  OnChanges, OnInit,
  Output,
  signal,
  SimpleChanges
} from '@angular/core';
import {JsonPipe, NgIf, NgStyle} from "@angular/common";
import {FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {ToastComponent} from "../../../../ui/toast/toast.component";
import {IRequirement, ISection} from "../../../../models/announcement";
import {Department, District, Province} from "../../../../models/ubigeos";

@Component({
  selector: 'app-section-form',
  standalone: true,
  imports: [
    JsonPipe,
    NgIf,
    ReactiveFormsModule,
    ToastComponent,
    NgStyle
  ],
  templateUrl: './section-form.component.html',
  styleUrl: './section-form.component.scss'
})
export class SectionFormComponent implements OnChanges {
  section = input.required<ISection>();
  formDynamic = input.required<FormGroup>();

  dataDepartment = input<Department[]>();
  dataProvince = input<Province[]>();
  dataDistrict = input<District[]>();

  signalDepartment = signal<Department[]>([]);
  signalProvince = signal<Province[]>([]);
  signalDistrict = signal<District[]>([]);

  recordRequirements = input.required<Record<string, string>>()

  public dataSource: Record<number, any> = {};

  protected isLoading: boolean = false;

  @Output() onProcessFile = new EventEmitter<{event: any, key: string, type: number}>();

  constructor() {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (this.dataDepartment()?.length) this.signalDepartment.set(this.dataDepartment() || []);
    this.loadDataSource();
  }

  public loadDataSource() {
    this.section().requisitos.forEach(req => {
      if (req.tipo_requisito_id === 4) {
        this.dataSource[req.id || 0] = this.getDataSourceByReq(req);
      }
    })
  }

  public getDataSourceByReq(req: IRequirement) {
    const type = req.opciones ?? '';
    let key = '';
    if (req.is_dependent) {
      key = this.getValueDependent(req.field_dependent!);
      if (!key) return [];
      debugger
      this.updateDataSource(type, key);
    }
    switch (type) {
      case 'department':
        return this.signalDepartment();
      case 'province':
        return this.signalProvince();
      case 'district':
        return this.signalDistrict();
      default:
        return this.loadOptionDefault(type);
    }
  }

  private updateDataSource(type: string, key: string) {
    switch (type) {
      case 'department':
        this.signalDepartment.set(this.dataDepartment() || []);
        break;
      case 'province':
        this.signalProvince.set(this.dataProvince()!.filter((item: any) => item.departament_id === key))
        break;
      case 'district':
        this.signalDistrict.set(this.dataDistrict()!.filter((item: any) => item.province_id === key))
        break;
      default:
    }
  }

  private getValueDependent(field: string): string {
    const reqTemp = this.recordRequirements()[field];
    return this.formDynamic().get(reqTemp.toString()!)?.value;
  }

  public loadOptionDefault(option: string): {id: string, name: string}[] {
    const optionArray = option.split('|');
    return optionArray.map((item: any) => ({ id: item, name: item }));
  }

  public processFile(event: Event, key: string, type: number) {
    this.onProcessFile.emit({event, key, type });
  }

  public isDependent(req: IRequirement): boolean {
    if (!req.is_dependent || !req.show_dependent) {
      return true;
    }

    const key = this.getValueDependent(req.field_dependent!);
    if (!key) {
      return false;
    }

    const control = this.formDynamic().get(req.id?.toString()!);
    if (key.toLowerCase() === req.value_dependent?.toLowerCase()) {
      control?.setValidators(Validators.required);
      control?.updateValueAndValidity();
      return true;
    } else {
      control?.clearValidators();
      control?.updateValueAndValidity();
      return false;
    }
  }

  public onCheckboxChange(event: any, key: string): void {
    let selectedValues = this.formDynamic().get(key)?.value || '';
    const checkboxValue = event.target.value;

    if (event.target.checked) {
      selectedValues = selectedValues
        ? `${selectedValues},${checkboxValue}`
        : checkboxValue;
    } else {
      const valuesArray = selectedValues.split(',').filter((value: string) => value !== checkboxValue);
      selectedValues = valuesArray.join(',');
    }

    this.formDynamic().get(key)?.setValue(selectedValues);
  }

  isChecked(value: string | number, key: string): boolean {
    const selectedValues = this.formDynamic().get(key)?.value || '';
    const valuesArray = selectedValues.split(',');
    return valuesArray.includes(value.toString());
  }
}
