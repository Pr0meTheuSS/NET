package nsu.ru.Lab3.controllers;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

class SearchTerm {
    // This should match the input field name in your form
    private String locationName; 

    public String getLocationName() {
        return locationName;
    }

    public void setLocationName(String locationName) {
        this.locationName = locationName;
    }
}

@Controller
public class AsyncApiRequestsController {
    @GetMapping("/search")
    public String searchPage(Model model) {
        // Create an instance of SearchTerm or set a value as needed
        model.addAttribute("searchTerm", new SearchTerm());
        return "welcome"; 
    }
}
