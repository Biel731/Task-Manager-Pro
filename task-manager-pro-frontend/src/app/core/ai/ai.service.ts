import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface SuggestTitleResponse {
  titles: string[];
}

export interface ImproveDescriptionResponse {
  improved_description: string;
  bullets: string[];
}

@Injectable({ providedIn: 'root' })
export class AiService {
  private baseUrl = '/api/ai';

  constructor(private http: HttpClient) {}

  suggestTitle(description: string): Observable<SuggestTitleResponse> {
    return this.http.post<SuggestTitleResponse>(`${this.baseUrl}/suggest-title`, {
      description,
    });
  }

  improveDescription(title: string, description: string): Observable<ImproveDescriptionResponse> {
    return this.http.post<ImproveDescriptionResponse>(`${this.baseUrl}/improve-description`, {
      title,
      description,
    });
  }
}
