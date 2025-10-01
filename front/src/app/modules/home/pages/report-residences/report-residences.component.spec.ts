import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReportResidencesComponent } from './report-residences.component';

describe('ReportResidencesComponent', () => {
  let component: ReportResidencesComponent;
  let fixture: ComponentFixture<ReportResidencesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReportResidencesComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ReportResidencesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
