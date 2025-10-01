import { Component } from '@angular/core';
import { Services } from "../../../../core/utils/statics/statics";
import { IServices } from "../../../../core/models/services";
import { NgIf } from '@angular/common';
import { Subscription } from 'rxjs';
import { ManagerService } from '../../../../core/services/manager/manager.service';

@Component({
    selector: 'app-services',
    standalone: true,
    imports: [NgIf],
    templateUrl: './services.component.html',
    styleUrl: './services.component.scss'
})
export class ServicesComponent {
    isActiveSRQ: boolean = false
    protected services: IServices[] = Services;
    private _subscriptions: Subscription = new Subscription();
    isModalOpen = false;

    public constructor(private _managerService: ManagerService) {
        this.validateActiveSRQ()
    }

    toggleModal(state: boolean): void {
        this.isModalOpen = state;
    }

    validateActiveSRQ() {
        this._subscriptions.add(
            this._managerService.hasActiveSRQ().subscribe({
                next: (response: any) => {
                    this.isActiveSRQ = response?.detalle?.activa || false;
                },
                error: (error: any) => {
                    console.error('Error al obtener las respuestas del estudiante:', error);
                    this.isActiveSRQ = false;
                }
            })
        );
    }

}
