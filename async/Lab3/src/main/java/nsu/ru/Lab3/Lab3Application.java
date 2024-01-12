package nsu.ru.Lab3;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import nsu.ru.Lab3.Configs.Configs;

@SpringBootApplication
public class Lab3Application {

	@Autowired
	Configs configs;

	public static void main(String[] args) {
		SpringApplication.run(Lab3Application.class, args);
	}

}
