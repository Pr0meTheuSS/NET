package nsu.ru.Lab3.controllers;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;


import nsu.ru.Lab3.LocationApi.*;
import nsu.ru.Lab3.PlacesApi.*;
import nsu.ru.Lab3.WeatherApi.*;

@Controller
public class MyController {
    // TODO: rename
    private LocationResponseDTO places = null;
    
    @Autowired
    private WeatherApiIface weatherApiService;

    @Autowired
    private LocationApiIface locationApiService;

    @Autowired
    private PlacesApiIface placesApiService;

    @GetMapping("/search")
    public String searchPage(Model model) {
        model.addAttribute("searchTerm", new SearchTerm());
        return "search"; 
    }

    @GetMapping("/locations")
    public String locations(@RequestParam String locationName, Model model) {
        try {
            places = locationApiService.fetchLocations(locationName);
            implaceLocationsIntoPage(places, model);
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }

        return "index";
    }


    @GetMapping("/info/{id}")
    public String myPageLocationsInfo(@PathVariable int id, Model model) {
        try {
            if (places == null) {
                return "error";
            }

            String lat = places.getHitsLat(id);
            String lon = places.getHitsLon(id);

            System.out.println("Before places in radius call");

            CompletableFuture<PlacesDTO> placesDataFuture = CompletableFuture.supplyAsync(() -> {
                try {
                    return placesApiService.fetchPlacesInRadius(lat, lon, "1000");
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });

            System.out.println("Before get description of places");
            CompletableFuture<List<PlaceInfoView>> placeInfoViewsFuture = placesDataFuture.thenApply(data -> {
                List<PlaceInfoView> placesInfoViews = new ArrayList<>();
                int limit = 0;
                for (Feature f : data.getFeatures()) {
                    if (limit++ > 10) {
                        break;
                    }

                    try {
                        PlaceInfo p = placesApiService.fetchPlaceDescriptionByXid(f.getId());
                        placesInfoViews.add(mapPlaceInfoTPlaceInfoView(p));
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }
                return placesInfoViews;
            });

            System.out.println("Before weather api call");
            CompletableFuture<WeatherData> weatherFuture = CompletableFuture.supplyAsync(() -> {
                try {
                    return weatherApiService.getWeatherAtPoint(lat, lon);
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });

            System.out.println("After weather api call");


            CompletableFuture<Void> allOf = CompletableFuture.allOf(weatherFuture, placesDataFuture, placeInfoViewsFuture);

            System.out.println("Before join");
            allOf.join();

            WeatherView wv = mapWeatherDTOtoView(weatherFuture.join());
            List<PlaceInfoView> placesInfoViews = placeInfoViewsFuture.join();

            System.out.println("After join");

            model.addAttribute("WeatherView", wv);
            model.addAttribute("placeInfoList", placesInfoViews);
        } catch (Exception e) {
            e.printStackTrace();
        }

        return "info";
    }

    private void implaceLocationsIntoPage(LocationResponseDTO dto, Model model) {
        List<PlaceView> itemList = new ArrayList<>();
        for (Location l: dto.getHits()) {
            String placeName = l.getCountry() + " " + l.getCity() + " " + l.getName();
            placeName = placeName.replaceAll("null", "");

            itemList.add(new PlaceView(dto.getHits().indexOf(l), placeName));
        }
        
        model.addAttribute("items", itemList);
    }

    private WeatherView mapWeatherDTOtoView(WeatherData data) {
        WeatherView ret = new WeatherView();

        ret.setDescription(data.getDescription());
        ret.setTemperature(String.valueOf(data.getTemp()));
        ret.setMinTemperature(String.valueOf(data.getTempMin()));
        ret.setMaxTemperature(String.valueOf(data.getTempMax()));
        ret.setWindSpeed(String.valueOf(data.getWindSpeed()));

        return ret;
    }

    private PlaceInfoView mapPlaceInfoTPlaceInfoView(PlaceInfo data) {
        PlaceInfoView ret = new PlaceInfoView();

        ret.setName(data.getName());
        ret.setAddress(data.getAddress());
        ret.setKinds(data.getKinds());

        return ret;
    }
}
