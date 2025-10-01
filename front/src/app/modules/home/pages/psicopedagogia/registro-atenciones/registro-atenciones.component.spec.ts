import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistroAtencionesComponent } from './registro-atenciones.component';

describe('RegistroAtencionesComponent', () => {
  let component: RegistroAtencionesComponent;
  let fixture: ComponentFixture<RegistroAtencionesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RegistroAtencionesComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RegistroAtencionesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
