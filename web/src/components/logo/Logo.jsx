import { Title } from "@mantine/core";
import styles from "./logo.module.css";

export default function Logo() {
  return (
    <Title size={28} textWrap="stable" className={styles.localFontText}>
      Gobali
    </Title>
  );
}
