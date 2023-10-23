package nsu.ru.Lab3.WeatherApi;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import com.fasterxml.jackson.databind.ObjectMapper;

public class WeatherApiImpl implements WeatherApiIface {
    private final HttpClient httpClient;
    public WeatherApiImpl() {
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public WeatherData getWeatherAtPoint(String lat, String lon) throws IOException, InterruptedException {
        try {
            String url = "https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&units=metric" +  "&lang=ru" + "&appid=" + "f64de4c7268561459c992f43532336c2";
            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .build();
            HttpResponse<String> resp;
            resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            System.out.println(resp.body());

            ObjectMapper objectMapper = new ObjectMapper();
            return objectMapper.readValue(resp.body(), WeatherData.class);
        } catch (IOException | InterruptedException e) {
           throw e;
        }    
    }    
}