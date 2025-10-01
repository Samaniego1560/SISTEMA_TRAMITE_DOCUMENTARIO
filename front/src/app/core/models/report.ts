export interface RequirementReport {
  name: string;
  type: number;
  value: string;
}

interface SectionReport {
  description: string;
  requirements: RequirementReport[];
}

interface AnnouncementReport {
  name: string | null;
  details_requests: SectionReport[];
}

export interface StudentFileReport {
  title: string;
  link_profile: string;
  announcements: AnnouncementReport[];
  first_section: SectionReport;
  year: string;
}
