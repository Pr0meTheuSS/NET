package nsu.ru.Lab3.WeatherApi;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

import nsu.ru.Lab3.Configs.Configs;

@Service
public class WeatherApiImpl implements WeatherApiIface {
    private final HttpClient httpClient;

    private final String getWeatherAtPointUrl = "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&units=metric&lang=ru&appid={apikey}";
    private String apikey;

    @Autowired
    public WeatherApiImpl(Configs cnfgs) {
        apikey = cnfgs.getWeatherApiKey();
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public WeatherData getWeatherAtPoint(String lat, String lon) throws IOException, InterruptedException {
        try {
            String url = getWeatherAtPointUrl;
            url = url.replace("{lat}", lat);
            url = url.replace("{lon}", lon);
            url = url.replace("{apikey}", apikey);
            System.out.println("Weather url: " + url);
            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .build();
            HttpResponse<String> resp;
            resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

            ObjectMapper objectMapper = new ObjectMapper();
            return objectMapper.readValue(resp.body(), WeatherData.class);
        } catch (IOException | InterruptedException e) {
           throw e;
        }    
    }    
}