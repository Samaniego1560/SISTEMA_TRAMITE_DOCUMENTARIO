import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistroPostAtencionComponent } from './registro-post-atencion.component';

describe('RegistroPostAtencionComponent', () => {
  let component: RegistroPostAtencionComponent;
  let fixture: ComponentFixture<RegistroPostAtencionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RegistroPostAtencionComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RegistroPostAtencionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
