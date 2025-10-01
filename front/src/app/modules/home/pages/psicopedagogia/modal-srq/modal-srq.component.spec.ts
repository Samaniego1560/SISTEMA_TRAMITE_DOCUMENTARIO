import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModalSrqComponent } from './modal-srq.component';

describe('ModalSrqComponent', () => {
  let component: ModalSrqComponent;
  let fixture: ComponentFixture<ModalSrqComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ModalSrqComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ModalSrqComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
