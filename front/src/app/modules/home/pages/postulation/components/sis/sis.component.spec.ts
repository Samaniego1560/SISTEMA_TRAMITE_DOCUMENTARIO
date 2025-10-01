import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SisComponent } from './sis.component';

describe('SisComponent', () => {
  let component: SisComponent;
  let fixture: ComponentFixture<SisComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SisComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SisComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
