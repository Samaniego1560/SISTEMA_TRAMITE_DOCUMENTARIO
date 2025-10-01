import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable,of} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class FileService {
  constructor(
    private http: HttpClient
  ) {
  }

  base64ToExcel(base64: string, filename: string) {
    console.log("Base64 recibido:", base64.slice(0, 50) + "..."); // Verificar los primeros caracteres
    try {
        // Convertir Base64 a Blob
        const byteCharacters = atob(base64);
        const byteArray = new Uint8Array(byteCharacters.length);

        for (let i = 0; i < byteCharacters.length; i++) {
            byteArray[i] = byteCharacters.charCodeAt(i);
        }

        const blob = new Blob([byteArray], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });

        console.log("Tamaño del archivo generado:", blob.size); // Depuración

        if (blob.size === 0) {
            throw new Error("El archivo generado está vacío o corrupto.");
        }

        // Crear un enlace para descargar el archivo
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);

        link.href = url;
        link.download =  `${filename}_${new Date().toISOString().slice(0, 10)}.xlsx`;

        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        URL.revokeObjectURL(url); // Liberar memoria
    } catch (error: any) {

    }

  }

}
