import {ACTION_TYPES} from "../../contans/areas/permissions.constant";
import {PermissionsService} from "../../services/permissions/permissions.service";
import {Directive, inject, Injectable, Input} from "@angular/core";

@Directive()
export abstract class FormBaseComponent<T> {
  @Input() infoData: T | undefined;
  private _permissionsService: PermissionsService = inject(PermissionsService)

  get isModeCreation() {
    return !this.infoData;
  }

  get canSaveForm() {
    return this.canCreate || this.canUpdate;
  }

  get canCreate() {
    return this.isModeCreation && this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.CREATE
    )
  }

  get canUpdate() {
    return !this.isModeCreation && this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.UPDATE
    )
  }

  abstract isValid(): boolean;
  abstract getData(): any;
  abstract get name(): string;
}

