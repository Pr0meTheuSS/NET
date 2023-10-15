package nsu.ru.Lab3.controllers;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.ArrayList;
import java.util.List;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

@Controller
public class MyController {

    private final HttpClient httpClient;

    public MyController() {
        this.httpClient = HttpClient.newHttpClient();
    }

    @GetMapping("/my-page")
    public String myPage(@RequestParam String locationName, Model model) {
        ObjectMapper objectMapper = new ObjectMapper();        
        String ret = fetchLocation(locationName);
        LocationResponseDTO dto = new LocationResponseDTO();
        try {
            dto =  objectMapper.readValue(ret, LocationResponseDTO.class);
            System.out.println(dto);
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }

        List<PlaceItem> itemList = new ArrayList<>();
        for (Location l: dto.getHits()) {

            String placeName = l.getCountry() + " " + l.getCity() + " " + l.getName();
            placeName = placeName.replaceAll("null", "");

            itemList.add(new PlaceItem(dto.getHits().indexOf(l), placeName));
        }
        model.addAttribute("items", itemList);

        return "index";
    }

    // TODO: вынести в асинхронный вызов.
    private String fetchLocation(String locationName) {
        String url = prepareUrlForfetchingLocation(locationName);
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();

        HttpResponse<String> resp;
        try {
            resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            return resp.body();
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }

        return "";
    }

    private String prepareUrlForfetchingLocation(String locationName) {
        String encodedValue = "";
        try {
            encodedValue = URLEncoder.encode(locationName, "UTF-8");
            System.out.println(encodedValue);
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }

        return "https://graphhopper.com/api/1/geocode?q=" + encodedValue + "&locale=ru&key=b21dcab5-cc27-472b-8a3e-d1eb62c38a04";
    }
}

class PlaceItem {
    private int placeId;
    private String placeName;

    public PlaceItem(int itemId, String placeName) {
        this.placeId = itemId;
        this.placeName = placeName;
    }

    public int getPlaceId() {
        return placeId;
    }

    public void setPlaceId(int itemId) {
        this.placeId = itemId;
    }

    public String getPlaceName() {
        return placeName;
    }

    public void setPlaceName(String placeName) {
        this.placeName = placeName;
    }
}
