package nsu.ru.Lab3.LocationApi;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.net.URLEncoder;

import com.fasterxml.jackson.databind.ObjectMapper;

import nsu.ru.Lab3.Configs.Configs;

@Service
public class LocationApiImpl implements LocationApiIface {
    private final HttpClient httpClient;
    private String fetchLocationsUrl = "https://graphhopper.com/api/1/geocode?q={locationName}&locale=ru&key={apikey}";
    private String apikey;

    @Autowired
    public LocationApiImpl(Configs cnfgs) {
        apikey = cnfgs.getLocationsApiKey();
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public LocationResponseDTO fetchLocations(String locationName) throws IOException, InterruptedException {
        ObjectMapper objectMapper = new ObjectMapper();
        String ret;
        try {
            ret = fetchLocation(locationName);
            return objectMapper.readValue(ret, LocationResponseDTO.class);
        } catch (IOException | InterruptedException e) {
            throw e;
        }
    }
    
    private String fetchLocation(String locationName) throws IOException, InterruptedException {
        String url = prepareUrlForfetchingLocation(locationName);
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();

        HttpResponse<String> resp;
        resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return resp.body();
    }

    private String prepareUrlForfetchingLocation(String locationName) throws UnsupportedEncodingException {
        String encodedLocationName = URLEncoder.encode(locationName, "UTF-8");
        String url = fetchLocationsUrl;
        url = url.replace("{locationName}", encodedLocationName);
        url = url.replace("{apikey}", apikey);
        return url;
    }
}
